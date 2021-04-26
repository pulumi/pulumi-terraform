package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v5/go/state/internal"
)

// HttpStateArgs specifies the configuration options for a Terraform Remote State
// stored in the HTTP backend.
type HttpStateArgs struct {
	// Address is the address of the HTTP endpoint.
	Address pulumi.StringInput

	// UpdateMethod is the HTTP method to use when updating state. Defaults to `POST`.
	UpdateMethod pulumi.StringPtrInput

	// LockAddress is the address of the lock REST endpoint. Not setting a value disables locking.
	LockAddress pulumi.StringPtrInput

	// LockMethod is the HTTP method to use when locking. Defaults to `LOCK`.
	LockMethod pulumi.StringPtrInput

	// UnlockAddress is the address of the unlock REST endpoint. Not setting a value disables locking.
	UnlockAddress pulumi.StringPtrInput

	// LockMethod is the HTTP method to use when unlocking. Defaults to `LOCK`.
	UnlockMethod pulumi.StringPtrInput

	// Username is the username used for HTTP basic authentication.
	Username pulumi.StringPtrInput

	// Password is the password used for HTTP basic authentication.
	Password pulumi.StringPtrInput

	// SkipCertVerification is whether to skip TLS verification. Defaults to `false`.
	SkipCertVerification pulumi.BoolPtrInput

	// Workspace is the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (l *HttpStateArgs) toInternalArgs() pulumi.Input {
	return internal.HttpStateReferenceArgs{
		BackendType:          pulumi.String("http"),
		Address:              l.Address,
		UpdateMethod:         l.UpdateMethod,
		LockAddress:          l.LockAddress,
		LockMethod:           l.LockMethod,
		UnlockAddress:        l.UnlockAddress,
		UnlockMethod:         l.UnlockMethod,
		Username:             l.Username,
		Password:             l.Password,
		SkipCertVerification: l.SkipCertVerification,
		Workspace:            l.Workspace,
	}
}

func (l *HttpStateArgs) validateArgs() error {
	if l.Address == pulumi.String("") {
		return errors.New("`Address` is a required parameter")
	}
	return nil
}
