package examples

import (
	"os"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
)

func getCwd(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.FailNow()
	}

	return cwd
}

func getBaseOptions() integration.ProgramTestOptions {
	return integration.ProgramTestOptions{
		RunUpdateTest: false,
		LocalProviders: []integration.LocalDependency{{
			Package: "terraform",
			Path:    "../bin",
		}},
	}
}

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
