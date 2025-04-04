// Copyright 2016-2025, Pulumi Corporation.  All rights reserved.
//go:build yaml || all
// +build yaml all

package examples

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAML(t *testing.T) {

	expectedStackOutputs := map[string]string{
		"local": `{
                "state": {
                    "bucket_arn": "arn:aws:s3:::hello-world-abc12345",
                    "public_subnet_ids": [
                        "subnet-023a5a6867d194162",
                        "subnet-0eea17cb6af21b5e5",
                        "subnet-02822dcd2e06634cf"
                    ],
                    "vpc_id": "vpc-0d9ff66ccda7c9765"
                }
            }`,
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
				Dir: testDir,
				ExportStateValidator: func(t *testing.T, state []byte) {
					if expected, ok := expectedStackOutputs[dir.Name()]; ok {
						var actual map[string]any
						require.NoError(t, json.Unmarshal(state, &actual))
						for _, r := range actual["deployment"].(map[string]any)["resources"].([]any) {
							res := r.(map[string]any)
							if res["type"].(string) == "pulumi:pulumi:Stack" {
								actualOutputs, err := json.Marshal(res["outputs"])
								require.NoError(t, err)
								assert.JSONEq(t, expected, string(actualOutputs))
								return
							}
						}
						assert.Fail(t, "Could not find a stack resource")
					}
				},
			})
			integration.ProgramTest(t, &opts)
		})
	}
}
