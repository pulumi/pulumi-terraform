package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LocalStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Path pulumi.StringInput
}

type localStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Path string `pulumi:"path"`
}

func (LocalStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*localStateReferenceArgs)(nil)).Elem()
}
