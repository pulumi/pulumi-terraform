// Copyright 2016-2026, Pulumi Corporation.
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
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// TestStateReferenceReadLocalNullOutput reads local state whose outputs contain a
// null value and checks that null survives to the RPC response instead of being
// conflated with "not present".
//
// The test goes through [p.RawServer] because the null/missing distinction is a
// property of the gRPC marshaling layer: before pulumi-go-provider v1.4.1 that
// layer dropped null-valued keys (SkipNulls), so null_field below vanished from
// the response.
func TestStateReferenceReadLocalNullOutput(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	// Terraform elides top-level null outputs from state, so the realistic way a
	// null reaches a state file is inside an object output.
	config := `
output "present" {
  value = "hello"
}

output "obj" {
  value = {
    set_field  = "yes"
    null_field = null
  }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(config), 0o600))
	tofu := func(args ...string) {
		cmd := exec.CommandContext(ctx, "tofu", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "tofu %v failed:\n%s", args, out)
	}
	tofu("init", "-input=false")
	tofu("apply", "-auto-approve", "-input=false")

	InitTfBackend()
	prov := infer.Provider(infer.Options{
		Functions: []infer.InferredFunction{infer.Function(&GetLocalReference{})},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{"state_reference": "state"},
	})
	server, err := p.RawServer("terraform", "6.0.0", prov)(nil)
	require.NoError(t, err)

	args, err := structpb.NewStruct(map[string]any{
		localPathAttribute: filepath.Join(dir, "terraform.tfstate"),
	})
	require.NoError(t, err)

	resp, err := server.Invoke(ctx, &pulumirpc.InvokeRequest{
		Tok:  "terraform:state:getLocalReference",
		Args: args,
	})
	require.NoError(t, err)
	require.Empty(t, resp.GetFailures())

	assert.Equal(t, map[string]any{
		"outputs": map[string]any{
			"present": "hello",
			"obj": map[string]any{
				"set_field":  "yes",
				"null_field": nil,
			},
		},
	}, resp.GetReturn().AsMap())
}
