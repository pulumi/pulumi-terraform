package shim

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-svchost/disco"
	be "github.com/hashicorp/terraform/internal/backend"
	backendInit "github.com/hashicorp/terraform/internal/backend/init"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InitTfBackend() { backendInit.Init(disco.New()) }

type RemoteStateReferenceInputs struct {
	// TODO: what is this for? Is it always be pulumi.String("remote")?
	BackendType string

	BackendConfig BackendConfig

	// Workspace is a struct specifying which remote workspace(s) to use.
	Workspaces WorkspaceStateArgs
}

type BackendConfig struct {
	// The name of the resource to read.
	ResourceName string

	// Organization is the name of the organization containing the targeted workspace(s).
	Organization string

	// Hostname is the remote backend hostname to which to connect. Defaults to `app.terraform.io`.
	Hostname string

	// Token is the token used to authenticate with the remote backend.
	Token string
}

// WorkspaceStateArgs specifies the configuration options for a workspace for use with the remote enhanced backend.
type WorkspaceStateArgs struct {
	// Name is the full name of one remote workspace. When configured, only the default workspace
	// can be used. This option conflicts with prefix.
	Name string

	// Prefix is the prefix used in the names of one or more remote workspaces, all of which can be used
	// with this configuration. If unset, only the default workspace can be used. This option
	// conflicts with name
	Prefix string
}

func RemoteStateReferenceRead(ctx context.Context, args RemoteStateReferenceInputs) (map[string]any, error) {
	backendType := args.BackendType

	// Ensure the backendType is known about by Terraform
	backendInitFn := backendInit.Backend(backendType)
	if backendInitFn == nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported backend type %q", backendType)
	}

	// If we have a workspace specified, get the value for that. Use the default otherwise
	workspaceName := be.DefaultStateName
	if args.Workspaces.Name != "" {
		workspaceName = args.Workspaces.Name
	}

	// Get the configuration schema from the backend
	backend := backendInitFn()

	// Attempt to coerce our config object into the config schema types - note errors
	backendConfigCoerced, err := backend.ConfigSchema().CoerceValue(cty.ObjectVal(map[string]cty.Value{
		"token":         cty.StringVal(args.BackendConfig.Token),
		"organization":  cty.StringVal(args.BackendConfig.Organization),
		"hostname":      cty.StringVal(args.BackendConfig.Hostname),
		"resource_name": cty.StringVal(args.BackendConfig.ResourceName),
	}))
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

	var outputs map[string]any
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

type LocalStateReferenceInputs struct {
	Path string
}

func LocalStateReferenceRead(ctx context.Context, args LocalStateReferenceInputs) (map[string]any, error) {
	backendType := "local"

	// Ensure the backendType is known about by Terraform
	backendInitFn := backendInit.Backend(backendType)
	if backendInitFn == nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported backend type %q", backendType)
	}

	// If we have a workspace specified, get the value for that. Use the default otherwise
	workspaceName := be.DefaultStateName
	// if args.Workspaces.Name != "" {
	// 	workspaceName = args.Workspaces.Name
	// }

	// Get the configuration schema from the backend
	backend := backendInitFn()

	// Attempt to coerce our config object into the config schema types - note errors
	backendConfigCoerced, err := backend.ConfigSchema().CoerceValue(cty.ObjectVal(map[string]cty.Value{
		"path": cty.StringVal(args.Path),
	}))
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
