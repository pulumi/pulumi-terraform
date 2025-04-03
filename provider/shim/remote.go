package shim

import (
	"context"

	"github.com/zclconf/go-cty/cty"
)

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
	return StateReferenceRead(ctx, "remote", args.Workspaces.Name, map[string]cty.Value{
		"token":         cty.StringVal(args.BackendConfig.Token),
		"organization":  cty.StringVal(args.BackendConfig.Organization),
		"hostname":      cty.StringVal(args.BackendConfig.Hostname),
		"resource_name": cty.StringVal(args.BackendConfig.ResourceName),
	})
}
