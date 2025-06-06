package shim

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-svchost/disco"
	backendInit "github.com/hashicorp/terraform/internal/backend/init"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InitTfBackend() { backendInit.Init(disco.New()) }

func StateReferenceRead(
	ctx context.Context,
	backendType string,
	workspaceName string,
	backendConfigValue map[string]cty.Value,
) (map[string]any, error) {
	// Ensure the backendType is known about by Terraform
	backendInitFn := backendInit.Backend(backendType)
	if backendInitFn == nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported backend type %q", backendType)
	}

	// Get the configuration schema from the backend
	backend := backendInitFn()

	// Attempt to coerce our config object into the config schema types - note errors
	backendConfigCoerced, err := backend.ConfigSchema().CoerceValue(cty.ObjectVal(backendConfigValue))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error coercing config from Pulumi format to cty: %s", err)
	}

	// Attempt to prepare the backend with configuration, returning any diagnostics to the engine
	preparedBackendConfig, diagnostics := backend.PrepareConfig(backendConfigCoerced)
	if diagnostics.HasErrors() {
		return nil, status.Errorf(codes.Internal, "error preparing config: %s", diagnostics.Err())
	}

	// Actually prepare the backend with the valid configuration
	diagnostics = backend.Configure(preparedBackendConfig)
	if diagnostics.HasErrors() {
		return nil, status.Errorf(codes.InvalidArgument, "error in backend configuration: %s",
			diagnostics.ErrWithWarnings())
	}

	// Get the state manager from the backend for the appropriate workspace
	stateManager, err := backend.StateMgr(workspaceName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error constructing backend state manager: %s", err)
	}

	// Refresh the state
	if err := stateManager.RefreshState(); err != nil {
		return nil, status.Errorf(codes.NotFound, "error refreshing Terraform state: %s", err)
	}

	// Check the state isn't empty
	state := stateManager.State()
	if state == nil {
		return nil, status.Error(codes.NotFound, "remote state not found")
	}

	// Convert back into the type that we expect.

	outputs := map[string]any{}
	for k, v := range state.RootModule().OutputValues {
		jsonBytes, err := ctyjson.Marshal(v.Value, v.Value.Type())
		if err != nil {
			return nil, fmt.Errorf("error marshaling cty to JSON: %w", err)
		}
		var goV any
		if err := json.Unmarshal(jsonBytes, &goV); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		outputs[k] = goV
	}
	return outputs, nil
}
