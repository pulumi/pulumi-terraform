package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v4/go/state/internal"
)

// EtcdV2StateArgs specifies the configuration options for a Terraform Remote State
// stored in the etcd v2 backend. Please not there is a separate func for the Etcd v2 backend
type EtcdV2StateArgs struct {
	// Path at which to store the state.
	Path pulumi.StringInput

	// Endpoints is a space-separated list of the etcd endpoints.
	Endpoints pulumi.StringInput

	// Username is the username with which to authenticate to etcd.
	Username pulumi.StringPtrInput

	// Password is the password with which to authenticate to etcd.
	Password pulumi.StringPtrInput

	// Workspace is the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (l *EtcdV2StateArgs) toInternalArgs() pulumi.Input {
	return internal.EtcdV2StateReferenceArgs{
		BackendType: pulumi.String("etcd"),
		Path:        l.Path,
		Endpoints:   l.Endpoints,
		Username:    l.Username,
		Password:    l.Password,
		Workspace:   l.Workspace,
	}
}

func (l *EtcdV2StateArgs) validateArgs() error {
	if l.Path == pulumi.String("") || l.Endpoints == pulumi.String("") {
		return errors.New("`Path` and `Endpoints` are required parameters")
	}
	return nil
}
