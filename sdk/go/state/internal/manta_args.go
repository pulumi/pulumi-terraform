package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type MantaStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Account               pulumi.StringInput
	User                  pulumi.StringPtrInput
	Url                   pulumi.StringPtrInput
	KeyMaterial           pulumi.StringPtrInput
	KeyID                 pulumi.StringInput
	Path                  pulumi.StringInput
	InsecureSkipTlsVerify pulumi.BoolInput
	Workspace             pulumi.StringPtrInput
}

type mantaStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Account               string  `pulumi:"account"`
	User                  *string `pulumi:"user"`
	Url                   *string `pulumi:"url"`
	KeyMaterial           *string `pulumi:"keyMaterial"`
	KeyID                 string  `pulumi:"keyId"`
	Path                  string  `pulumi:"path"`
	InsecureSkipTlsVerify bool    `pulumi:"insecureSkipTlsVerify"`
	Workspace             *string `pulumi:"workspace"`
}

func (MantaStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*mantaStateReferenceArgs)(nil)).Elem()
}
