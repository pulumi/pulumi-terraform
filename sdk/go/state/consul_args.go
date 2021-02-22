package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v4/go/state/internal"
)

// ConsulArgs specifies the configuration options for a Terraform Remote State
// stored in the Consul backend
type ConsulArgs struct {
	// Path is the path in the Consul KV store
	Path pulumi.StringInput

	// AccessToken is the Consul Access Token. Sourced from `CONSUL_HTTP_TOKEN` in the environment, if unset.
	AccessToken pulumi.StringInput

	// Address is the DNS name and port of the Consul HTTP endpoint specified in the format `dnsname:port`. Defaults
	// to the local agent HTTP listener.
	Address pulumi.StringPtrInput

	// Scheme specifies which protocol to use when talking to the given address - either `http` or `https`. TLS
	// support can also be enabled by setting the environment variable `CONSUL_HTTP_SSL` to `true`.
	Schema pulumi.StringPtrInput

	// Datacenter is the datacenter to use. Defaults to that of the agent.
	Datacenter pulumi.StringPtrInput

	// HttpAuth is the HTTP Basic Authentication credentials to be used when communicating with Consul, in the
	// format of either `user` or `user:pass`. Sourced from `CONSUL_HTTP_AUTH`, if unset.
	HttpAuth pulumi.StringPtrInput

	// Gzip is whether to compress the state data using gzip. Set to `true` to compress the state data using gzip,
	// or `false` (default) to leave it uncompressed.
	Gzip pulumi.BoolPtrInput

	// CAFile is the path to a PEM-encoded certificate authority used to verify the remote agent's certificate.
	// Sourced from `CONSUL_CAFILE` in the environment, if unset.
	CAFile pulumi.StringPtrInput

	// CertFile is the path to a PEM-encoded certificate provided to the remote agent; requires use of key_file.
	// Sourced from `CONSUL_CLIENT_CERT` in the environment, if unset.
	CertFile pulumi.StringPtrInput

	// KeyFile is the path to a PEM-encoded private key, required if cert_file is specified. Sourced from
	// `CONSUL_CLIENT_KEY` in the environment, if unset.
	KeyFile pulumi.StringPtrInput

	// Workspace os the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (a *ConsulArgs) toInternalArgs() pulumi.Input {
	return internal.ConsulStateReferenceArgs{
		BackendType: pulumi.String("consul"),
		Path:        a.Path,
		AccessToken: a.AccessToken,
		Address:     a.Address,
		Scheme:      a.Schema,
		Datacenter:  a.Datacenter,
		HttpAuth:    a.HttpAuth,
		Gzip:        a.Gzip,
		CAFile:      a.CAFile,
		CertFile:    a.CertFile,
		KeyFile:     a.KeyFile,
		Workspace:   a.Workspace,
	}
}

func (l *ConsulArgs) validateArgs() error {
	if l.Path == pulumi.String("") || l.AccessToken == pulumi.String("") {
		return errors.New("`Path` and `Path` are required parameters")
	}
	return nil
}
