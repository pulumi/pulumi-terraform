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
	"testing"

	"github.com/hashicorp/terraform/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func ptr[T any](v T) *T { return &v }

func TestWorkspacesStateMgrName(t *testing.T) {
	tests := []struct {
		name     string
		ws       Workspaces
		expected string
	}{
		{
			name:     "name set passes default to StateMgr",
			ws:       Workspaces{Name: ptr("my-workspace")},
			expected: defaultWorkspace,
		},
		{
			name:     "prefix set passes prefix to StateMgr",
			ws:       Workspaces{Prefix: ptr("app-")},
			expected: "app-",
		},
		{
			name:     "both nil passes empty string",
			ws:       Workspaces{},
			expected: "",
		},
		{
			name:     "name takes precedence over prefix",
			ws:       Workspaces{Name: ptr("my-workspace"), Prefix: ptr("app-")},
			expected: defaultWorkspace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.ws.stateMgrName())
		})
	}
}

func TestStateReferenceReadLocal(t *testing.T) {
	InitTfBackend()

	outputs, err := shim.StateReferenceRead(
		context.Background(), "local", defaultWorkspace, map[string]cty.Value{
			"path": cty.StringVal("testdata/test.tfstate"),
		},
	)
	require.NoError(t, err)

	assert.Equal(t, "hello", outputs["greeting"])
	assert.Equal(t, float64(42), outputs["count"])
}

func TestStateReferenceReadUnsupportedBackend(t *testing.T) {
	InitTfBackend()

	_, err := shim.StateReferenceRead(
		context.Background(), "nonexistent", defaultWorkspace, map[string]cty.Value{},
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported backend type")
}

// TestRemoteBackendIgnoresVersionConflict verifies that StateReferenceRead
// calls IgnoreVersionConflict on backends that support it, allowing reads
// even when the local TF version doesn't match the workspace's version.
// We test this indirectly: the remote backend is instantiated and we confirm
// it implements the interface that StateReferenceRead uses.
func TestRemoteBackendIgnoresVersionConflict(t *testing.T) {
	InitTfBackend()

	// Get the remote backend the same way the shim does
	backendInitFn := shim.BackendFactory("remote")
	require.NotNil(t, backendInitFn, "remote backend should be registered")

	backend := backendInitFn()
	_, ok := backend.(shim.VersionConflictIgnorer)
	assert.True(t, ok, "remote backend should implement IgnoreVersionConflict()")
}
