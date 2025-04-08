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
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRemoteBackendOrganization(t *testing.T) string {
	org, found := os.LookupEnv("TFE_ORGANIZATION")
	if !found {
		t.Skipf("Skipping... cannot find TFE_ORGANIZATION")
	}

	return org
}

func getRemoteBackendToken(t *testing.T) string {
	token, found := os.LookupEnv("TFE_TOKEN")
	if !found {
		t.Skipf("Skipping... cannot find TFE_TOKEN")
	}

	return token
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

func LanguageTests(t *testing.T, language string) {
	type languageTest struct {
		doesNotNeedConfig    bool
		expectedStackOutputs map[string]any
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
			if !test.doesNotNeedConfig {
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
