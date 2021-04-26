package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v5/go/state/internal"
)

// RemoteBackendStateArgs specifies the configuration options for a Terraform Remote State
// stored in the remote backend.
type RemoteBackendStateArgs struct {
	// Organization is the name of the organization containing the targeted workspace(s).
	Organization pulumi.StringInput

	// Hostname is the remote backend hostname to which to connect. Defaults to `app.terraform.io`.
	Hostname pulumi.StringPtrInput

	// Token is the token used to authenticate with the remote backend.
	Token pulumi.StringPtrInput

	// Workspace is a struct specifying which remote workspace(s) to use.
	Workspaces WorkspaceStateArgs
}

// WorkspaceStateArgs specifies the configuration options for a workspace for use with the remote enhanced backend.
type WorkspaceStateArgs struct {
	// Name is the full name of one remote workspace. When configured, only the default workspace
	// can be used. This option conflicts with prefix.
	Name pulumi.StringPtrInput

	// Prefix is the prefix used in the names of one or more remote workspaces, all of which can be used
	// with this configuration. If unset, only the default workspace can be used. This option
	// conflicts with name
	Prefix pulumi.StringPtrInput
}

func (l *RemoteBackendStateArgs) toInternalArgs() pulumi.Input {
	args := internal.RemoteBackendStateReferenceArgs{
		BackendType:  pulumi.String("remote"),
		Organization: l.Organization,
		Token:        l.Token,
		Hostname:     l.Hostname,
	}

	internalArgs := internal.WorkspaceArgs{}
	if l.Workspaces.Name != nil {
		internalArgs.Name = l.Workspaces.Name
	}
	if l.Workspaces.Prefix != nil {
		internalArgs.Prefix = l.Workspaces.Prefix
	}

	args.Workspaces = internalArgs

	return args
}

func (l *RemoteBackendStateArgs) validateArgs() error {
	if l.Organization == pulumi.String("") {
		return errors.New("`Organization` is a required parameter")
	}
	return nil
}
