// Copyright 2016-2025, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package examples

import (
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func getRemoteBackendOrganization(t *testing.T) string { return getEnv(t, "TFE_ORGANIZATION") }

func getRemoteBackendToken(t *testing.T) string { return getEnv(t, "TFE_TOKEN") }

func getEnv(t *testing.T, env string) string {
	value, found := os.LookupEnv(env)
	if !found {
		if os.Getenv("CI") != "" {
			t.Fatalf("Failing... cannot find %s", env)
		}
		t.Skipf("Skipping... cannot find %s", env)
	}
	return value
}

func getwd(t *testing.T) string {
	s, err := os.Getwd()
	require.NoError(t, err)
	return s
}

func getDependencies(t *testing.T, language string) []string {
	switch language {
	case "go":
		return []string{"github.com/pulumi/pulumi-terraform/sdk/v6=" + getwd(t) + "/../sdk"}
	case "nodejs":
		return []string{"@pulumi/terraform"}
	case "python":
		return []string{"../sdk/python"}
	case "dotnet":
		return []string{"Pulumi.Terraform"}
	default:
		return nil
	}
}

func TestMain(m *testing.M) {
	set := func(envVar, path string) error {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		return os.Setenv(envVar, abs)
	}

	err := errors.Join(
		set("PULUMI_LOCAL_NUGET", "../nuget"),
		set("PULUMI_LOCAL_MAVEN", "../maven"),
	)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// createTFCWorkspace creates a temporary TFC workspace with the given Terraform version
// and uploads a state file with the given outputs. It registers a cleanup function to
// delete the workspace when the test finishes. Returns the workspace name.
func createTFCWorkspace(t *testing.T, org, token, tfVersion string, outputs map[string]any) string {
	t.Helper()
	ctx := context.Background()

	client, err := tfe.NewClient(&tfe.Config{Token: token})
	require.NoError(t, err, "creating TFC client")

	wsName := fmt.Sprintf("pulumi-tf-inttest-%d", time.Now().UnixNano())

	ws, err := client.Workspaces.Create(ctx, org, tfe.WorkspaceCreateOptions{
		Name:             tfe.String(wsName),
		TerraformVersion: tfe.String(tfVersion),
		ExecutionMode:    tfe.String("local"),
	})
	require.NoError(t, err, "creating TFC workspace")

	t.Cleanup(func() {
		_ = client.Workspaces.DeleteByID(context.Background(), ws.ID)
	})

	// Build a minimal tfstate with the desired outputs.
	type tfOutput struct {
		Value any    `json:"value"`
		Type  string `json:"type"`
	}
	tfOutputs := map[string]tfOutput{}
	for k, v := range outputs {
		tfOutputs[k] = tfOutput{Value: v, Type: "string"}
	}
	stateFile := map[string]any{
		"version":           4,
		"terraform_version": tfVersion,
		"serial":            1,
		"lineage":           "integration-test",
		"outputs":           tfOutputs,
		"resources":         []any{},
	}

	stateJSON, err := json.Marshal(stateFile)
	require.NoError(t, err, "marshaling tfstate")

	h := crypto.MD5.New() //nolint:gosec // MD5 is required by the TFC State Versions API
	h.Write(stateJSON)
	md5Str := fmt.Sprintf("%x", h.Sum(nil))
	stateB64 := base64.StdEncoding.EncodeToString(stateJSON)
	serial := int64(1)

	// The TFC API requires the workspace to be locked before uploading state.
	_, err = client.Workspaces.Lock(ctx, ws.ID, tfe.WorkspaceLockOptions{
		Reason: tfe.String("uploading integration test state"),
	})
	require.NoError(t, err, "locking workspace")

	_, err = client.StateVersions.Create(ctx, ws.ID, tfe.StateVersionCreateOptions{
		MD5:    tfe.String(md5Str),
		Serial: &serial,
		State:  tfe.String(stateB64),
	})
	require.NoError(t, err, "uploading state version")

	_, err = client.Workspaces.Unlock(ctx, ws.ID)
	require.NoError(t, err, "unlocking workspace")

	return wsName
}

func LanguageTests(t *testing.T, language string) {
	type languageTest struct {
		doesNotNeedConfig    bool
		expectedStackOutputs map[string]any
		// configFunc, when set, overrides the default config generation.
		configFunc func(t *testing.T) map[string]string
	}
	tests := map[string]languageTest{
		"local": {
			doesNotNeedConfig: true,
			expectedStackOutputs: map[string]any{
				"state": map[string]any{
					"bucket_arn": "arn:aws:s3:::hello-world-abc12345",
					"public_subnet_ids": []any{
						"subnet-023a5a6867d194162",
						"subnet-0eea17cb6af21b5e5",
						"subnet-02822dcd2e06634cf",
					},
					"vpc_id": "vpc-0d9ff66ccda7c9765",
				},
				"bucketArn":     "arn:aws:s3:::hello-world-abc12345",
				"firstSubnetId": "subnet-023a5a6867d194162",
			},
		},
		"remote": {
			expectedStackOutputs: map[string]any{
				"state": map[string]any{
					"4dabf18193072939515e22adb298388d": "1b47061264138c4ac30d75fd1eb44270",
					"plaintext":                        "{\"password\":\"EOZcr9x4V@ep8T1gjmR4RJ39aT9vQDsDwZx\"}",
				},
			},
		},
		"remote-name": {
			// The output is wrapped in Pulumi's secret envelope because the token
			// input is marked as secret (provider:"secret" on GetRemoteReferenceArgs.Token),
			// and secret-ness propagates to outputs. WireDependencies/NeverSecret is
			// intended to prevent this but is not yet implemented for infer-based
			// functions (https://github.com/pulumi/pulumi-go-provider/issues/323).
			// With DecryptSecretsInOutput the plaintext is visible, but still wrapped:
			//   "4dabf18193072939515e22adb298388d" = secret sentinel key
			//   "1b47061264138c4ac30d75fd1eb44270" = secret signature value
			//   "plaintext" = the actual outputs serialized as a JSON string
			expectedStackOutputs: map[string]any{
				"state": map[string]any{
					"4dabf18193072939515e22adb298388d": "1b47061264138c4ac30d75fd1eb44270",
					"plaintext":                        `{"test_output":"hello-from-integration-test"}`,
				},
			},
			configFunc: func(t *testing.T) map[string]string {
				token := getRemoteBackendToken(t)
				org := getRemoteBackendOrganization(t)
				wsName := createTFCWorkspace(t, org, token, "1.9.8", map[string]any{
					"test_output": "hello-from-integration-test",
				})
				return map[string]string{
					"remote_tf_token": token,
					"remote_tf_org":   org,
					"workspace_name":  wsName,
				}
			},
		},
	}

	examplesDir := "."
	examples, err := os.ReadDir(examplesDir)
	require.NoError(t, err)

	for _, dir := range examples {
		if !dir.IsDir() {
			continue
		}

		testDir := filepath.Join(examplesDir, dir.Name(), language)
		_, err := os.Stat(testDir)
		if os.IsNotExist(err) {
			continue
		}
		require.NoError(t, err)

		// We have found a yaml example called dir.Name, so we should run it async.
		t.Run(dir.Name(), func(t *testing.T) {
			test := tests[dir.Name()]
			var config map[string]string
			if test.configFunc != nil {
				config = test.configFunc(t)
			} else if !test.doesNotNeedConfig {
				config = map[string]string{
					"remote_tf_token": getRemoteBackendToken(t),
					"remote_tf_org":   getRemoteBackendOrganization(t),
				}
			}
			opts := integration.ProgramTestOptions{
				Dir:                    testDir,
				DecryptSecretsInOutput: true,
				Config:                 config,
				LocalProviders: []integration.LocalDependency{{
					Package: "terraform",
					Path:    "../bin",
				}},
				Dependencies: getDependencies(t, language),
				ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
					assert.Equal(t, test.expectedStackOutputs, stack.Outputs)
				},
			}
			integration.ProgramTest(t, &opts)
		})
	}
}
