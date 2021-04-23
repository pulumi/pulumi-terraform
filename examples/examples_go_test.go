// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.
// +build go all

package examples

import (
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func TestGoLocal013(t *testing.T) {
	t.Skip("temp skipping while preping for major version change")
	test := getGoBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-go"),
			Config: map[string]string{
				"statefile": "terraform.0-13-0.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestGoLocal012(t *testing.T) {
	t.Skip("temp skipping while preping for major version change")
	test := getGoBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-go"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestGoS3013(t *testing.T) {
	t.Skip("temp skipping while preping for major version change")
	test := getGoBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-go"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-11-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestGoS3012(t *testing.T) {
	t.Skip("temp skipping while preping for major version change")
	test := getGoBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-go"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-12-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestGoRemoteBackend(t *testing.T) {
	t.Skip("temp skipping while preping for major version change")
	test := getGoBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "remote-backend-go"),
			Config: map[string]string{
				"organization":  getRemoteBackendOrganization(t),
				"workspaceName": "dev",
			},
			Secrets: map[string]string{
				"tfeToken": getRemoteBackendToken(t),
			},
		})

	integration.ProgramTest(t, &test)
}

func getGoBaseOptions(t *testing.T) integration.ProgramTestOptions {
	t.Skip("temp skipping while preping for major version change")
	base := getBaseOptions()
	baseGo := base.With(integration.ProgramTestOptions{
		RunUpdateTest: false,
		Dependencies: []string{
			"github.com/pulumi/pulumi-terraform/sdk/v3/go",
		},
	})

	return baseGo
}
