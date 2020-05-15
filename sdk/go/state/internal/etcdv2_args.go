package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type EtcdV2StateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Path      pulumi.StringInput
	Endpoints pulumi.StringInput
	Username  pulumi.StringPtrInput
	Password  pulumi.StringPtrInput
	Workspace pulumi.StringPtrInput
}

type etcdV2StateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Path      string  `pulumi:"path"`
	Endpoints string  `pulumi:"endpoints"`
	Username  *string `pulumi:"username"`
	Password  *string `pulumi:"password"`
	Workspace *string `pulumi:"workspace"`
}

func (EtcdV2StateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*etcdV2StateReferenceArgs)(nil)).Elem()
}
