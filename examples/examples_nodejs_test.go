// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.
// +build nodejs all

package examples

import (
	"path"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func TestJSLocal013(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-nodejs"),
			Config: map[string]string{
				"statefile": "terraform.0-13-0.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSLocal012(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-nodejs"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSS3013(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-nodejs"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-13-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSS3012(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-nodejs"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-12-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSOss013(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "ossstate-nodejs"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"prefix":     "0-13-state",
				"region":     "us-west-1",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSOss012(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "ossstate-nodejs"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"prefix":     "0-12-state",
				"region":     "us-west-1",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSRemoteBackend(t *testing.T) {
	test := getJSBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "remote-backend-nodejs"),
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

func getJSBaseOptions(t *testing.T) integration.ProgramTestOptions {
	base := getBaseOptions()
	baseJS := base.With(integration.ProgramTestOptions{
		Dependencies: []string{
			"@pulumi/terraform",
		},
	})

	return baseJS
}
