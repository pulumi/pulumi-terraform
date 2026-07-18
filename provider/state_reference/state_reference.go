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
	"github.com/zclconf/go-cty/cty"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// defaultWorkspace is the sentinel name Terraform backends use for the
// implicit, always-present workspace.
const defaultWorkspace = "default"

func InitTfBackend() { shim.InitTfBackend() }

type StateReferenceOutputs struct {
	// Outputs is a map of the outputs from the Terraform state file
	Outputs map[string]any `pulumi:"outputs"`
}

var _ = (infer.Annotated)((*StateReferenceOutputs)(nil))

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

func ctyBoolOrNil(v *bool) cty.Value {
	if v == nil {
		return cty.NullVal(cty.Bool)
	}
	return cty.BoolVal(*v)
}

func ctyIntOrNil(v *int) cty.Value {
	if v == nil {
		return cty.NullVal(cty.Number)
	}
	return cty.NumberIntVal(int64(*v))
}

func ctyStringSetOrNil(v []string) cty.Value {
	if len(v) == 0 {
		return cty.NullVal(cty.Set(cty.String))
	}
	elems := make([]cty.Value, len(v))
	for i, s := range v {
		elems[i] = cty.StringVal(s)
	}
	return cty.SetVal(elems)
}

func ctyStringMapOrNil(v map[string]string) cty.Value {
	if len(v) == 0 {
		return cty.NullVal(cty.Map(cty.String))
	}
	elems := make(map[string]cty.Value, len(v))
	for k, s := range v {
		elems[k] = cty.StringVal(s)
	}
	return cty.MapVal(elems)
}

func stringOrZero(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
