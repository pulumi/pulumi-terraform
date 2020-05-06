package examples

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
)

func TestJSLocal011(t *testing.T) {
	test := getJSBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-nodejs"),
			Config: map[string]string{
				"statefile": "terraform.0-11-3.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSLocal012(t *testing.T) {
	test := getJSBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-nodejs"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSS3011(t *testing.T) {
	test := getJSBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-nodejs"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-11-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestJSS3012(t *testing.T) {
	test := getJSBaseOptions().
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

func TestPyLocal011(t *testing.T) {
	test := getPyBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-python"),
			Config: map[string]string{
				"statefile": "terraform.0-11-3.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestPyLocal012(t *testing.T) {
	test := getPyBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-python"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestPyS3011(t *testing.T) {
	test := getPyBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-python"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-11-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestPyS3012(t *testing.T) {
	test := getPyBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-python"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-12-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetLocal011(t *testing.T) {
	test := getDotNetBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-dotnet"),
			Config: map[string]string{
				"statefile": "terraform.0-11-3.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetLocal012(t *testing.T) {
	test := getDotNetBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-dotnet"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetS3011(t *testing.T) {
	test := getDotNetBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "s3state-dotnet"),
			Config: map[string]string{
				"bucketName": "pulumi-terraform-remote-state-testing",
				"key":        "0-11-state",
				"region":     "us-west-2",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestDotNetS3012(t *testing.T) {
	test := getDotNetBaseOptions().
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

func TestGoLocal011(t *testing.T) {
	test := getGoBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-go"),
			Config: map[string]string{
				"statefile": "terraform.0-11-3.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestGoLocal012(t *testing.T) {
	test := getGoBaseOptions().
		With(integration.ProgramTestOptions{
			Dir: path.Join(getCwd(t), "localstate-go"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
		})

	integration.ProgramTest(t, &test)
}

func TestGoS3011(t *testing.T) {
	test := getGoBaseOptions().
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
	test := getGoBaseOptions().
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

func getCwd(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.FailNow()
	}

	return cwd
}

func getJSBaseOptions() integration.ProgramTestOptions {
	return integration.ProgramTestOptions{
		RunUpdateTest: false,
		Quick:         true,
		Dependencies: []string{
			"@pulumi/terraform",
		},
	}
}

func getPyBaseOptions() integration.ProgramTestOptions {
	return integration.ProgramTestOptions{
		RunUpdateTest: false,
		Quick:         true,
		Dependencies: []string{
			filepath.Join("..", "sdk", "python", "bin"),
		},
	}
}

func getDotNetBaseOptions() integration.ProgramTestOptions {
	return integration.ProgramTestOptions{
		RunUpdateTest: false,
		Quick:         true,
		Dependencies: []string{
			"Pulumi.Terraform",
		},
	}
}

func getGoBaseOptions() integration.ProgramTestOptions {
	return integration.ProgramTestOptions{
		RunUpdateTest: false,
		Quick:         true,
		Dependencies: []string{
			"github.com/pulumi/pulumi-terraform/sdk/v2",
		},
	}
}
