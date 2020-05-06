package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type GcsStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Bucket        pulumi.StringInput
	Credentials   pulumi.StringPtrInput
	Prefix        pulumi.StringPtrInput
	EncryptionKey pulumi.StringPtrInput
	Workspace     pulumi.StringPtrInput
}

type gcsStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Bucket        string  `pulumi:"bucket"`
	Credentials   *string `pulumi:"credentials"`
	Prefix        *string `pulumi:"prefix"`
	EncryptionKey *string `pulumi:"encryptionKey"`
	Workspace     *string `pulumi:"workspace"`
}

func (GcsStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*gcsStateReferenceArgs)(nil)).Elem()
}
