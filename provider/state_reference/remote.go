// Copyright 2016-2025, Pulumi Corporation.
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
	"github.com/zclconf/go-cty/cty"
)

type RemoteStateReference struct{}

var (
	_ = (infer.Annotated)((*RemoteStateReference)(nil))
	_ = (infer.ExplicitDependencies[RemoteStateReferenceInputs, StateReferenceOutputs])((*RemoteStateReference)(nil))
)

func (r *RemoteStateReference) Annotate(a infer.Annotator) {
	a.Describe(&r, "Access state from a remote backend.")
}

// Taken from https://developer.hashicorp.com/terraform/language/backend/remote#configuration-variables
type RemoteStateReferenceInputs struct {
	Hostname     *string    `pulumi:"hostname,optional"`
	Organization string     `pulumi:"organization"`
	Token        *string    `pulumi:"token,optional" provider:"secret"`
	Workspaces   Workspaces `pulumi:"workspaces"`
}

func (r *RemoteStateReferenceInputs) Annotate(a infer.Annotator) {
	a.Describe(&r.Hostname, "The remote backend hostname to connect to.")
	a.Describe(&r.Organization, "The name of the organization containing the targeted workspace(s).")
	a.Describe(&r.Token, "The token used to authenticate with the remote backend.")

	a.SetDefault(&r.Hostname, "app.terraform.io")
}

// WireDependencies lets us tell user's that our outputs shouldn't be secret, even when
// the token (when provided) is always secret.
//
// TODO[https://github.com/pulumi/pulumi-go-provider/issues/323]: This doesn't currently
// work; [infer.ExplicitDependencies] is not currently implemented for [infer] based
// functions.
func (r *RemoteStateReference) WireDependencies(
	f infer.FieldSelector, _ *RemoteStateReferenceInputs, state *StateReferenceOutputs,
) {
	f.OutputField(&state).NeverSecret() // The output should never be secret by default
}

type Workspaces struct {
	Name   *string `pulumi:"name,optional"`
	Prefix *string `pulumi:"prefix,optional"`
}

func (r *Workspaces) Annotate(a infer.Annotator) {
	a.Describe(&r.Name, "The full name of one remote workspace. When configured, only the default workspace can be "+
		"used. This option conflicts with prefix.")
	a.Describe(&r.Prefix, "A prefix used in the names of one or more remote workspaces, all of which can be used "+
		"with this configuration. The full workspace names are used in HCP Terraform, and the short names "+
		"(minus the prefix) are used on the command line for Terraform CLI workspaces. If omitted, only "+
		"the default workspace can be used. This option conflicts with name.")
}

func (r *RemoteStateReference) Call(
	ctx context.Context, args RemoteStateReferenceInputs,
) (StateReferenceOutputs, error) {
	results, err := shim.StateReferenceRead(ctx, "remote", stringOrZero(args.Workspaces.Name), map[string]cty.Value{
		"hostname":     ctyStringOrNil(args.Hostname),
		"organization": cty.StringVal(args.Organization),
		"token":        ctyStringOrNil(args.Token),
		"workspaces": cty.ObjectVal(map[string]cty.Value{
			"name":   ctyStringOrNil(args.Workspaces.Name),
			"prefix": ctyStringOrNil(args.Workspaces.Prefix),
		}),
	})

	return StateReferenceOutputs{results}, err
}
