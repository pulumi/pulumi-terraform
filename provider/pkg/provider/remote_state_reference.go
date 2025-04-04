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

var _ = (infer.Annotated)((*RemoteStateReference)(nil))

func (r *RemoteStateReference) Annotate(a infer.Annotator) {
	a.Describe(&r, "TODO")
}

// Taken from https://developer.hashicorp.com/terraform/language/backend/remote#configuration-variables
type RemoteStateReferenceInputs struct {
	Hostname     *string    `pulumi:"hostname,optional"`
	Organization string     `pulumi:"organization"`
	Token        *string    `pulumi:"token,optional"`
	Workspaces   Workspaces `pulumi:"workspaces"`
}

type Workspaces struct {
	Name   *string `pulumi:"name,optional"`
	Prefix *string `pulumi:"prefix,optional"`
}

func (r *RemoteStateReference) Call(
	ctx context.Context, args RemoteStateReferenceInputs,
) (StateReferenceOutputs, error) {
	results, err := shim.StateReferenceRead(ctx, "remote", "", map[string]cty.Value{
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
