// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi-terraform/v6/provider/state_reference"
	"github.com/pulumi/pulumi-terraform/v6/provider/version"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

const (
	Name = "terraform"
)

// This provider uses the `pulumi-go-provider` library to produce a code-first provider definition.
func NewProvider() p.Provider {
	pkg := infer.Provider(infer.Options{
		// This is the metadata for the provider
		Metadata: schema.Metadata{
			DisplayName: "Terraform",
			Description: "The Terraform provider for Pulumi lets you consume the outputs " +
				"contained in Terraform state from your Pulumi programs.",
			Keywords: []string{
				"terraform",
				"kind/native",
				"category/utility",
			},
			Homepage:   "https://pulumi.com",
			License:    "Apache-2.0",
			Repository: "https://github.com/pulumi/pulumi-terraform",
			Publisher:  "Pulumi",
			LogoURL:    "https://raw.githubusercontent.com/pulumi/pulumi-terraform-provider/main/assets/logo.png",
			LanguageMap: map[string]any{
				"csharp": map[string]any{
					"respectSchemaVersion": true,
					"packageReferences": map[string]string{
						"Pulumi": "3.*",
					},
				},
				"go": map[string]any{
					"respectSchemaVersion":           true,
					"generateResourceContainerTypes": true,
					"importBasePath": fmt.Sprintf(
						"github.com/pulumi/pulumi-terraform/sdk/v%d/go/terraform", version.Version.Major),
				},
				"nodejs": map[string]any{
					"respectSchemaVersion": true,
				},
				"python": map[string]any{
					"respectSchemaVersion": true,
					"pyproject": map[string]bool{
						"enabled": true,
					},
				},
				"java": map[string]any{
					"buildFiles":                      "gradle",
					"gradleNexusPublishPluginVersion": "2.0.0",
					"dependencies": map[string]any{
						"com.pulumi:pulumi":               "1.0.0",
						"com.google.code.gson:gson":       "2.8.9",
						"com.google.code.findbugs:jsr305": "3.0.2",
					},
				},
			},
		},
		Functions: []infer.InferredFunction{
			infer.Function[*provider.GetLocalReference](),
			infer.Function[*provider.GetRemoteReference](),
		},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"state_reference": "state",
		},
	})

	{
		// Initialize the TF back-end exactly once during provider configuration
		oldConfigure := pkg.Configure
		pkg.Configure = func(ctx context.Context, req p.ConfigureRequest) error {
			NewTerraformLogRedirector(ctx)
			provider.InitTfBackend()
			if oldConfigure != nil {
				return oldConfigure(ctx, req)
			}
			return nil
		}
	}

	return pkg
}
