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
	a.Describe(&r, "TODO")
}

type LocalStateReferenceInputs struct {
	Path string `pulumi:"path"`
}

func (r *LocalStateReference) Call(
	ctx context.Context, args LocalStateReferenceInputs,
) (StateReferenceOutputs, error) {
	results, err := shim.StateReferenceRead(ctx, "local", "", map[string]cty.Value{
		"path": cty.StringVal(args.Path),
	})

	return StateReferenceOutputs{results}, err
}
