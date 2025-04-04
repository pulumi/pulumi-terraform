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

type LocalStateReference struct{}

var _ = (infer.Annotated)((*LocalStateReference)(nil))

func (r *LocalStateReference) Annotate(a infer.Annotator) {
	a.Describe(&r, "Access state from the local filesystem.")
}

// Taken from https://developer.hashicorp.com/terraform/language/backend/local#configuration-variables
type LocalStateReferenceInputs struct {
	Path         *string `pulumi:"path,optional"`
	WorkspaceDir *string `pulumi:"workspaceDir,optional"`
}

func (r *LocalStateReferenceInputs) Annotate(a infer.Annotator) {
	a.Describe(&r.Path, `The path to the tfstate file. This defaults to "terraform.tfstate" relative to the root module by default.`)
	a.Describe(&r.WorkspaceDir, `The path to non-default workspaces.`)
}

func (r *LocalStateReference) Call(
	ctx context.Context, args LocalStateReferenceInputs,
) (StateReferenceOutputs, error) {
	results, err := shim.StateReferenceRead(ctx, "local", "", map[string]cty.Value{
		"path":          ctyStringOrNil(args.Path),
		"workspace_dir": ctyStringOrNil(args.WorkspaceDir),
	})

	return StateReferenceOutputs{results}, err
}
