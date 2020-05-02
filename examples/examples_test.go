package examples

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
)

func TestExamples(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err, "expected a valid working directory: %v", err) {
		return
	}

	// base options shared amongst all tests.
	base := integration.ProgramTestOptions{
		RunUpdateTest:    false,
		SkipExportImport: true,
	}

	baseJS := base.With(integration.ProgramTestOptions{
		Dependencies: []string{
			"@pulumi/terraform",
		},
	})

	basePython := base.With(integration.ProgramTestOptions{
		Dependencies: []string{
			filepath.Join("..", "sdk", "python", "bin"),
		},
	})

	baseDotNet := base.With(integration.ProgramTestOptions{
		Dependencies: []string{
			"Pulumi.Terraform",
		},
	})

	tests := []integration.ProgramTestOptions{

		baseJS.With(integration.ProgramTestOptions{
			StackName: "js-tf0-11-3",
			Dir:       path.Join(cwd, "localstate-nodejs"),
			Config: map[string]string{
				"statefile": "terraform.0-11-3.tfstate",
			},
			NoParallel: true,
		}),
		baseJS.With(integration.ProgramTestOptions{
			StackName: "js-tf0-12-24",
			Dir:       path.Join(cwd, "localstate-nodejs"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
			NoParallel: true,
		}),
		basePython.With(integration.ProgramTestOptions{
			StackName: "py-tf0-11-3",
			Dir:       path.Join(cwd, "localstate-python"),
			Config: map[string]string{
				"statefile": "terraform.0-11-3.tfstate",
			},
			NoParallel: true,
		}),
		basePython.With(integration.ProgramTestOptions{
			StackName: "py-tf0-12-24",
			Dir:       path.Join(cwd, "localstate-python"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
			NoParallel: true,
		}),

		baseDotNet.With(integration.ProgramTestOptions{
			StackName: "dotnet-tf0-11-3",
			Dir:       path.Join(cwd, "localstate-dotnet"),
			Config: map[string]string{
				"statefile": "terraform.0-11-3.tfstate",
			},
			NoParallel: true,
		}),

		baseDotNet.With(integration.ProgramTestOptions{
			StackName: "dotnet-tf0-12-24",
			Dir:       path.Join(cwd, "localstate-dotnet"),
			Config: map[string]string{
				"statefile": "terraform.0-12-24.tfstate",
			},
			NoParallel: true,
		}),
	}

	for _, ex := range tests {
		example := ex
		t.Run(example.Dir, func(t *testing.T) {
			t.Log(example.StackName)
			integration.ProgramTest(t, &example)
		})
	}
}
