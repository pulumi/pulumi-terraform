package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v5/go/state/internal"
)

// LocalStateArgs specifies the configuration options for a Terraform Remote State
// stored in the local enhanced backend
type LocalStateArgs struct {
	// Path to the Terraform state file
	Path pulumi.StringInput
}

func (l *LocalStateArgs) toInternalArgs() pulumi.Input {
	return internal.LocalStateReferenceArgs{
		BackendType: pulumi.String("local"),
		Path:        l.Path,
	}
}

func (l *LocalStateArgs) validateArgs() error {
	if l.Path == pulumi.String("") {
		return errors.New("`Path` is a required parameter")
	}
	return nil
}
