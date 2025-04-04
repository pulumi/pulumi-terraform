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
	"github.com/hashicorp/terraform/shim"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/zclconf/go-cty/cty"
)

func InitTfBackend() { shim.InitTfBackend() }

// These are the outputs (or properties) of a LocalStateReference resource.
type StateReferenceOutputs struct {
	// Outputs is a map of the outputs from the Terraform state file
	Outputs map[string]any `pulumi:"outputs"`
}

var _ = (infer.Annotated)((*LocalStateReference)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (r *StateReferenceOutputs) Annotate(a infer.Annotator) {
	a.Describe(&r, "The result of fetching from a Terraform state store.")
	a.Describe(&r.Outputs, "The outputs displayed from Terraform state.")
}

func ctyStringOrNil(v *string) cty.Value {
	if v == nil {
		return cty.NullVal(cty.String)
	}
	return cty.StringVal(*v)
}
