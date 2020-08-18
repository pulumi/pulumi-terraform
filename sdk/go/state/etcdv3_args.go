package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v3/go/state/internal"
)

// EtcdV3StateArgs specifies the configuration options for a Terraform Remote State
// stored in the etcd v3 backend. Please not there is a separate func for the Etcd v3 backend
type EtcdV3StateArgs struct {
	// Path at which to store the state.
	Path pulumi.StringInput

	// Endpoints is a list of the etcd endpoints.
	Endpoints pulumi.StringArrayInput

	// Username is the username with which to authenticate to etcd. Sourced from `ETCDV3_USERNAME` env var if unset.
	Username pulumi.StringPtrInput

	// Password is the password with which to authenticate to etcd. Sourced from `ETCDV3_PASSWORD` env var if unset.
	Password pulumi.StringPtrInput

	// Prefix is an optional prefix to be added to keys when storing state in etcd.
	Prefix pulumi.StringPtrInput

	// CACertPath is the path to a PEM-encoded certificate authority bundle with which to verify certificates
	// of TLS-enabled etcd servers.
	CACertPath pulumi.StringPtrInput

	// CertPath is the path to a PEM-encoded certificate to provide to etcd for client authentication.
	CertPath pulumi.StringPtrInput

	// KeyPath is the path to a PEM-encoded key to use for client authentication.
	KeyPath pulumi.StringPtrInput

	// Workspace is the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (l *EtcdV3StateArgs) toInternalArgs() pulumi.Input {
	return internal.EtcdV3StateReferenceArgs{
		BackendType: pulumi.String("etcdv3"),
		Path:        l.Path,
		Endpoints:   l.Endpoints,
		Username:    l.Username,
		Password:    l.Password,
		Prefix:      l.Prefix,
		CertPath:    l.CertPath,
		CACertPath:  l.CACertPath,
		KeyPath:     l.KeyPath,
		Workspace:   l.Workspace,
	}
}

func (l *EtcdV3StateArgs) validateArgs() error {
	if l.Path == pulumi.String("") || l.Endpoints == nil {
		return errors.New("`Path` and `Endpoints` are required parameters")
	}
	return nil
}
