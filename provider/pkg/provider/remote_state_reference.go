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

// This is the type that implements the RemoteStateReference resource methods.
// The methods are declared in the read_resource.go file.
type RemoteStateReference struct{}

// The following statement is not required. It is a type assertion to indicate to Go that RemoteStateReference
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*RemoteStateReference)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (r *RemoteStateReference) Annotate(a infer.Annotator) {
	a.Describe(&r, "")
}

type RemoteStateReferenceInputs struct {
	// TODO: what is this for? Is it always be pulumi.String("remote")?
	BackendType string `pulumi:"backendType"`

	BackendConfig BackendConfig `pulumi:"backendConfig"`

	// Workspace is a struct specifying which remote workspace(s) to use.
	Workspaces Workspace `pulumi:"workspaces"`
}

type BackendConfig struct {
	// The name of the resource to read.
	ResourceName string `pulumi:"resourceName,optional"`

	// Organization is the name of the organization containing the targeted workspace(s).
	Organization string `pulumi:"organization"`

	// Hostname is the remote backend hostname to which to connect. Defaults to `app.terraform.io`.
	Hostname string `pulumi:"hostname,optional"`

	// Token is the token used to authenticate with the remote backend.
	Token string `pulumi:"token"`
}

// Workspace specifies the configuration options for a workspace for use with the remote enhanced backend.
type Workspace struct {
	// Name is the full name of one remote workspace. When configured, only the default workspace
	// can be used. This option conflicts with prefix.
	Name string `pulumi:"name,optional"`

	// Prefix is the prefix used in the names of one or more remote workspaces, all of which can be used
	// with this configuration. If unset, only the default workspace can be used. This option
	// conflicts with name
	Prefix string `pulumi:"prefix,optional"`
}

// These are the outputs (or properties) of a RemoteStateReference resource.
type RemoteStateReferenceOutputs struct {
	// Outputs is a map of the outputs from the Terraform state file
	Outputs map[string]any `pulumi:"outputs"`
}

func InitTfBackend() { shim.InitTfBackend() }

// Call implements the infer.Fn interface for RemoteStateReference.
func (r RemoteStateReference) Call(
	ctx context.Context, inputs RemoteStateReferenceInputs,
) (RemoteStateReferenceOutputs, error) {
	// Implement the logic for the Call method here.
	// Replace the following line with actual implementation.
	results, err := shim.RemoteStateReferenceRead(ctx, shim.RemoteStateReferenceInputs{
		BackendType: inputs.BackendType,
		BackendConfig: shim.BackendConfig{
			ResourceName: inputs.BackendConfig.ResourceName,
			Organization: inputs.BackendConfig.Organization,
			Hostname:     inputs.BackendConfig.Hostname,
			Token:        inputs.BackendConfig.Token,
		},
		Workspaces: shim.WorkspaceStateArgs{
			Name:   inputs.Workspaces.Name,
			Prefix: inputs.Workspaces.Prefix,
		},
	})
	if err != nil {
		return RemoteStateReferenceOutputs{}, err
	}
	return RemoteStateReferenceOutputs{results}, nil
}
