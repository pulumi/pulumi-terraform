package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v3/go/state/internal"
)

// PostgresStateArgs specifies the configuration options for a Terraform Remote State
// stored in the Postgres backend.
type PostgresStateArgs struct {
	// Bucket is the Postgres connection string; a `postgres://` URL.
	ConnStr pulumi.StringInput

	// SchemaName is the name of the automatically-managed Postgres schema. Defaults to `terraform_remote_state`.
	SchemaName pulumi.StringPtrInput

	// Workspace is the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (l *PostgresStateArgs) toInternalArgs() pulumi.Input {
	return internal.PostgresStateReferenceArgs{
		BackendType: pulumi.String("pg"),
		ConnStr:     l.ConnStr,
		SchemaName:  l.SchemaName,
		Workspace:   l.Workspace,
	}
}

func (l *PostgresStateArgs) validateArgs() error {
	if l.ConnStr == pulumi.String("") {
		return errors.New("`ConnStr` is a required parameter")
	}
	return nil
}
