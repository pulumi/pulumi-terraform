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
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// This is the type that implements the Run function methods.
// The methods are declared in the runController.go file.
type ReadResource struct{}

// Implementing Annotate lets you provide descriptions and default values for functions and they will
// be visible in the provider's schema and the generated SDKs.
func (r *ReadResource) Annotate(a infer.Annotator) {
	a.Describe(&r, "A local command to be executed.\n"+
		"This command will always be run on any preview or deployment. "+
		"Use `local.Command` to avoid duplicating executions.")
}

type ReadResourceInputs struct {
	// TODO: what is this for? Is it always be pulumi.String("remote")?
	BackendType pulumi.StringPtrInput

	// The name of the resource to read.
	ResourceName pulumi.StringInput

	// Organization is the name of the organization containing the targeted workspace(s).
	Organization pulumi.StringInput

	// Hostname is the remote backend hostname to which to connect. Defaults to `app.terraform.io`.
	Hostname pulumi.StringPtrInput

	// Token is the token used to authenticate with the remote backend.
	Token pulumi.StringPtrInput

	// Workspace is a struct specifying which remote workspace(s) to use.
	Workspaces WorkspaceStateArgs
}

// WorkspaceStateArgs specifies the configuration options for a workspace for use with the remote enhanced backend.
type WorkspaceStateArgs struct {
	// Name is the full name of one remote workspace. When configured, only the default workspace
	// can be used. This option conflicts with prefix.
	Name pulumi.StringPtrInput

	// Prefix is the prefix used in the names of one or more remote workspaces, all of which can be used
	// with this configuration. If unset, only the default workspace can be used. This option
	// conflicts with name
	Prefix pulumi.StringPtrInput
}

// Implementing Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (r *ReadResourceInputs) Annotate(a infer.Annotator) {
	a.Describe(&r.ResourceName, "The resource to read.")
}

type ReadResourceOutputs struct {
	ReadResourceInputs

	// TODO: do we need anything else in the output? Because ReadResource() reads the resource, so
	// the output should be on the resource itself.
}

// This is the Call method. It takes a ReadResourceInputs parameter and reads the resource specified in
// it.
func (*ReadResource) Call(ctx *pulumi.Context, input ReadResourceInputs) (ReadResourceOutputs, error) {
	output := ReadResourceOutputs{ReadResourceInputs: input}
	err := ctx.ReadResource(
		"terraform:state:RemoteStateReference",
		input.ResourceName,
		pulumi.ID(input.ResourceName),
		input,
		&output,
		// TODO: support options (pulumi.ResourceOption)
		// opts...,
	)
	return output, err
}
