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
	_ "embed"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed terraform.md
var resourceDoc string

// This is the type that implements the RemoteStateReference resource methods.
// The methods are declared in the read_resource.go file.
type RemoteStateReference struct{}

// The following statement is not required. It is a type assertion to indicate to Go that RemoteStateReference
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*RemoteStateReference)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *RemoteStateReference) Annotate(a infer.Annotator) {
	a.Describe(&c, resourceDoc)
}

// These are the inputs (or arguments) to a RemoteStateReference resource.
type RemoteStateReferenceInputs struct {
	ResourceInputs
}

// These are the outputs (or properties) of a RemoteStateReference resource.
type RemoteStateReferenceOutputs struct {
	// TODO: why do we need to include the input here in the output?
	RemoteStateReferenceInputs

	pulumi.CustomResourceState

	// Outputs is a map of the outputs from the Terraform state file
	Outputs pulumi.MapOutput `pulumi:"outputs"`
}
