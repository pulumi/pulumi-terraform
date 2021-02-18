package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v4/go/state/internal"
)

// MantaStateArgs specifies the configuration options for a Terraform Remote State
// stored in the Postgres backend.
type MantaStateArgs struct {
	// Account is the name of the Manta account. Sourced from `SDC_ACCOUNT` or `_ACCOUNT` in the environment, if unset.
	Account pulumi.StringInput

	// User is the username of the Manta account with which to authenticate.
	User pulumi.StringPtrInput

	// Url is the Manta API Endpoint. Sourced from `MANTA_URL` in the environment, if unset.
	// Defaults to `https://us-east.manta.joyent.com`.
	Url pulumi.StringPtrInput

	// KeyMaterial is the private key material corresponding with the SSH key whose fingerprint is
	// specified in keyId. Sourced from `SDC_KEY_MATERIAL` or `TRITON_KEY_MATERIAL`
	// in the environment, if unset. If no value is specified, the local SSH agent
	// is used for signing requests.
	KeyMaterial pulumi.StringPtrInput

	// KeyID is the fingerprint of the public key matching the key material specified in
	// keyMaterial, or in the local SSH agent.
	KeyID pulumi.StringInput

	// Path is the path relative to your private storage directory (`/$MANTA_USER/stor`)
	// where the state file will be stored.
	Path pulumi.StringInput

	// InsecureSkipTlsVerify is whether to skip verifying the TLS certificate presented by the Manta endpoint.
	// This can be useful for installations which do not have a certificate signed by a trusted
	// root CA. Defaults to false.
	InsecureSkipTlsVerify pulumi.BoolInput

	// Workspace is the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (l *MantaStateArgs) toInternalArgs() pulumi.Input {
	return internal.MantaStateReferenceArgs{
		BackendType:           pulumi.String("manta"),
		Account:               l.Account,
		User:                  l.User,
		Url:                   l.Url,
		KeyMaterial:           l.KeyMaterial,
		KeyID:                 l.KeyID,
		Path:                  l.Path,
		InsecureSkipTlsVerify: l.InsecureSkipTlsVerify,
		Workspace:             l.Workspace,
	}
}

func (l *MantaStateArgs) validateArgs() error {
	if l.Account == pulumi.String("") || l.KeyID == pulumi.String("") ||
		l.Path == pulumi.String("") || l.InsecureSkipTlsVerify == nil {
		return errors.New("`Account`, `KeyID`, `Path` and `InsecureSkipTlsVerify` are required parameters")
	}
	return nil
}
