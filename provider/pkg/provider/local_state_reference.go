// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"

	"github.com/hashicorp/terraform/shim"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// This is the type that implements the LocalStateReference resource methods.
// The methods are declared in the read_resource.go file.
type LocalStateReference struct{}

// The following statement is not required. It is a type assertion to indicate to Go that LocalStateReference
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*LocalStateReference)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (r *LocalStateReference) Annotate(a infer.Annotator) {
	a.Describe(&r, "")
}

type LocalStateReferenceInputs struct {
	Path string `pulumi:"path"`
}

// These are the outputs (or properties) of a LocalStateReference resource.
type LocalStateReferenceOutputs struct {
	// Outputs is a map of the outputs from the Terraform state file
	Outputs map[string]any `pulumi:"outputs"`
}

// Call implements the infer.Fn interface for LocalStateReference.
func (r *LocalStateReference) Call(
	ctx context.Context, inputs LocalStateReferenceInputs,
) (LocalStateReferenceOutputs, error) {
	// Implement the logic for the Call method here.
	// Replace the following line with actual implementation.
	results, err := shim.LocalStateReferenceRead(ctx, shim.LocalStateReferenceInputs{
		Path: inputs.Path,
	})
	if err != nil {
		return LocalStateReferenceOutputs{}, err
	}
	return LocalStateReferenceOutputs{results}, nil
}
