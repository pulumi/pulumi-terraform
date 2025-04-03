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

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi-terraform/provider_new/pkg/provider"
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
			Description: "TODO",
			Keywords: []string{
				"pulumi",
				"terraform",
				// TODO: where to find the keywords/tags?
				// "category/utility",
				// "kind/native",
			},
			Homepage:   "https://pulumi.com",
			License:    "Apache-2.0",
			Repository: "https://github.com/pulumi/pulumi-terraform",
			Publisher:  "Pulumi",
			LogoURL:    "TODO",
			// This contains language specific details for generating the provider's SDKs
			LanguageMap: map[string]any{
				// TODO: are these the same for all providers?
				"csharp": map[string]any{
					"respectSchemaVersion": true,
					"packageReferences": map[string]string{
						"Pulumi": "3.*",
					},
				},
				"go": map[string]any{
					"respectSchemaVersion":           true,
					"generateResourceContainerTypes": true,
					"importBasePath":                 "github.com/pulumi/pulumi-terraform/sdk/go/state",
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
		// Functions or invokes that are provided by the provider.
		Functions: []infer.InferredFunction{
			// The Read function is commented extensively for new pulumi-go-provider developers.
			infer.Function[provider.RemoteStateReference, provider.RemoteStateReferenceInputs, provider.RemoteStateReferenceOutputs](),
		},
	})

	{ // Initialize the TF back-end exactly once during provider configuration
		oldConfigure := pkg.Configure
		pkg.Configure = func(ctx context.Context, req p.ConfigureRequest) error {
			provider.InitTfBackend()
			return oldConfigure(ctx, req)
		}
	}

	return pkg
}
