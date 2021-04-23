package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v5/go/state/internal"
)

// GcsStateArgs specifies the configuration options for a Terraform Remote State
// stored in the Google Cloud Storage backend.
type GcsStateArgs struct {
	// Bucket is the name of the Google Cloud Storage bucket.
	Bucket pulumi.StringInput

	// Credentials is the path to Google Cloud Platform account credentials in JSON format. Sourced from
	// `GOOGLE_CREDENTIALS` in the environment if unset. If no value is provided Google
	// Application Default Credentials are used.
	Credentials pulumi.StringPtrInput

	// Prefix is the prefix used inside the Google Cloud Storage bucket. Named states for workspaces
	// are stored in an object named `&lt;prefix&gt;/&lt;name&gt;.tfstate`.
	Prefix pulumi.StringPtrInput

	// EncryptionKey is a 32 byte, base64-encoded customer supplied encryption key used to encrypt the
	// state. Sourced from `GOOGLE_ENCRYPTION_KEY` in the environment, if unset.
	EncryptionKey pulumi.StringPtrInput

	// Workspace is the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (l *GcsStateArgs) toInternalArgs() pulumi.Input {
	return internal.GcsStateReferenceArgs{
		BackendType:   pulumi.String("gcs"),
		Bucket:        l.Bucket,
		Credentials:   l.Credentials,
		Prefix:        l.Prefix,
		EncryptionKey: l.EncryptionKey,
		Workspace:     l.Workspace,
	}
}

func (l *GcsStateArgs) validateArgs() error {
	if l.Bucket == pulumi.String("") {
		return errors.New("`Bucket` is a required parameter")
	}
	return nil
}
