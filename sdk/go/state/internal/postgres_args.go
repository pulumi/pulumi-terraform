package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type PostgresStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	ConnStr    pulumi.StringInput
	SchemaName pulumi.StringPtrInput
	Workspace  pulumi.StringPtrInput
}

type postgresStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	ConnStr    string  `pulumi:"bucket"`
	SchemaName *string `pulumi:"schemaName"`
	Workspace  *string `pulumi:"workspace"`
}

func (PostgresStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*postgresStateReferenceArgs)(nil)).Elem()
}
