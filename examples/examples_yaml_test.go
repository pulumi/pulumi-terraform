// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build yaml || all
// +build yaml all

package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAML(t *testing.T) {

	expectedStackOutputs := map[string]any{
		"local": map[string]any{
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
		"remote": map[string]any{
			"name": map[string]any{},
			"prefix": map[string]any{
				"4dabf18193072939515e22adb298388d": "1b47061264138c4ac30d75fd1eb44270",
				"plaintext":                        "{\"password\":\"EOZcr9x4V@ep8T1gjmR4RJ39aT9vQDsDwZx\"}",
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

		testDir := filepath.Join(examplesDir, dir.Name(), "yaml")
		_, err := os.Stat(testDir)
		if os.IsNotExist(err) {
			continue
		}
		require.NoError(t, err)

		// We have found a yaml example called dir.Name, so we should run it async.
		t.Run(dir.Name(), func(t *testing.T) {
			opts := getBaseOptions().With(integration.ProgramTestOptions{
				Dir:                    testDir,
				DecryptSecretsInOutput: true,
				Config: map[string]string{
					"remote_tf_token": getRemoteBackendToken(t),
					"remote_tf_org":   getRemoteBackendOrganization(t),
				},
				ExtraRuntimeValidation: func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
					assert.Equal(t, expectedStackOutputs[dir.Name()], stack.Outputs)
				},
			})
			integration.ProgramTest(t, &opts)
		})
	}
}
