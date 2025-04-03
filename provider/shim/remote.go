package shim

import (
	"context"

	be "github.com/hashicorp/terraform/internal/backend"
	"github.com/zclconf/go-cty/cty"
)

type RemoteStateReferenceInputs struct {
	BackendConfig BackendConfig

	// Workspace is a struct specifying which remote workspace(s) to use.
	Workspaces WorkspaceStateArgs
}

type BackendConfig struct {
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
	// If we have a workspace specified, get the value for that. Use the default otherwise
	if args.Workspaces.Name == "" {
		args.Workspaces.Name = be.DefaultStateName
	}

	return StateReferenceRead(ctx, "remote", args.Workspaces.Name, map[string]cty.Value{
		"token":        cty.StringVal(args.BackendConfig.Token),
		"organization": cty.StringVal(args.BackendConfig.Organization),
		"hostname":     cty.StringVal(args.BackendConfig.Hostname),
		"workspaces": cty.ObjectVal(map[string]cty.Value{
			"name":   cty.StringVal(args.Workspaces.Name),
			"prefix": cty.StringVal(args.Workspaces.Prefix),
		}),
	})
}
