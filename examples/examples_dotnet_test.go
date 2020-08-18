// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.
// +build dotnet all

package examples

import (
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
)

func TestDotNetLocal013(t *testing.T) {
	test := getDotNetBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-dotnet"),
			Config: map[string]string{
				"statefile": "terraform.0-13-0.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetLocal012(t *testing.T) {
	test := getDotNetBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-dotnet"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetS3013(t *testing.T) {
	test := getDotNetBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-dotnet"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-13-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetS3012(t *testing.T) {
	test := getDotNetBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-dotnet"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-12-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetRemoteBackend(t *testing.T) {
	test := getDotNetBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "remote-backend-dotnet"),
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

func getDotNetBaseOptions(t *testing.T) integration.ProgramTestOptions {
	base := getBaseOptions()
	baseCsharp := base.With(integration.ProgramTestOptions{
		Dependencies: []string{
			"Pulumi.Terraform",
		},
	})

	return baseCsharp
}
