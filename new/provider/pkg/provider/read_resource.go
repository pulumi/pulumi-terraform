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

import "context"

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
	"github.com/pulumi/pulumi-go-provider/infer"
)

// This is the type that implements the Run function methods.
// The methods are declared in the runController.go file.
type ReadResource struct{}

// Implementing Annotate lets you provide descriptions and default values for functions and they will
// be visible in the provider's schema and the generated SDKs.
func (r *Run) Annotate(a infer.Annotator) {
	a.Describe(&r, "A local command to be executed.\n"+
		"This command will always be run on any preview or deployment. "+
		"Use `local.Command` to avoid duplicating executions.")
}

type ReadResourceInputs struct {
	BaseInputs
	Command string `pulumi:"command"`
}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (r *ReadResourceInputs) Annotate(a infer.Annotator) {
	a.Describe(&r.Command, "The resource to read.")
}

type ReadResourceOutputs struct {
	ReadResourceInputs
	BaseOutputs
}

// This is the Call method. It takes a ReadResourceInputs parameter and reads the resource specified in
// it.
func (*ReadResource) Call(ctx context.Context, input ReadResourceInputs) (ReadResourceOutputs, error) {
	r := ReadResourceOutputs{RunInputs: input}
	err := read(ctx, input.Command, r.RunInputs.BaseInputs, &r.BaseOutputs, input.Logging)
	return r, err
}
