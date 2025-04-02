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

	"github.com/pulumi/pulumi-command/provider/pkg/provider/common"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed command.md
var resourceDoc string

// This is the type that implements the Command resource methods.
// The methods are declared in the commandController.go file.
type RemoteStateReference struct{}

// The following statement is not required. It is a type assertion to indicate to Go that Command
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
	common.ResourceInputs

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

// These are the outputs (or properties) of a RemoteStateReference resource.
type RemoteStateReferenceOutputs struct {
	RemoteStateReferenceInputs
	// BaseOutputs
}
