package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type EtcdV3StateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Path       pulumi.StringInput
	Endpoints  pulumi.StringArrayInput
	Username   pulumi.StringPtrInput
	Password   pulumi.StringPtrInput
	Prefix     pulumi.StringPtrInput
	CACertPath pulumi.StringPtrInput
	CertPath   pulumi.StringPtrInput
	KeyPath    pulumi.StringPtrInput
	Workspace  pulumi.StringPtrInput
}

type etcdV3StateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Path       string  `pulumi:"path"`
	Endpoints  string  `pulumi:"endpoints"`
	Username   *string `pulumi:"username"`
	Password   *string `pulumi:"password"`
	Prefix     *string `pulumi:"prefix"`
	CACertPath *string `pulumi:"cacertPath"`
	CertPath   *string `pulumi:"certPath"`
	KeyPath    *string `pulumi:"keyPath"`
	Workspace  *string `pulumi:"workspace"`
}

func (EtcdV3StateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*etcdV3StateReferenceArgs)(nil)).Elem()
}
