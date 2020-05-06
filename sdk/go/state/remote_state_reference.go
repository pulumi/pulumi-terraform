package state

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// RemoteStateReference is a resource which allows reading the output from local or remote Terraform state files
type RemoteStateReference struct {
	pulumi.CustomResourceState

	// Outputs is a map of the outputs from the Terraform state file
	Outputs pulumi.MapOutput `pulumi:"outputs"`
}

// RemoteStateReferenceArgs is an interface implemented by all valid backend type Arg structs
type RemoteStateReferenceArgs interface {
	toInternalArgs() pulumi.Input
	validateArgs() error
}

// NewRemoteStateReference manages a reference to a Terraform Remote State.
// The root outputs of the remote state are available via the `Outputs` property
func NewRemoteStateReference(ctx *pulumi.Context, name string, args RemoteStateReferenceArgs,
	opts ...pulumi.ResourceOption) (*RemoteStateReference, error) {

	if args == nil {
		return nil, errors.New("missing required argument 'args'")
	}

	err := args.validateArgs()
	if err != nil {
		return nil, err
	}

	internalArgs := args.toInternalArgs()

	var resource RemoteStateReference
	err = ctx.ReadResource("terraform:state:RemoteStateReference", name, pulumi.ID(name), internalArgs, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}
