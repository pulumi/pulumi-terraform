package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type HttpStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Address              pulumi.StringInput
	UpdateMethod         pulumi.StringPtrInput
	LockAddress          pulumi.StringPtrInput
	LockMethod           pulumi.StringPtrInput
	UnlockAddress        pulumi.StringPtrInput
	UnlockMethod         pulumi.StringPtrInput
	Username             pulumi.StringPtrInput
	Password             pulumi.StringPtrInput
	SkipCertVerification pulumi.BoolPtrInput
	Workspace            pulumi.StringPtrInput
}

type httpStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Address              string  `pulumi:"address"`
	UpdateMethod         *string `pulumi:"updateMethod"`
	LockAddress          *string `pulumi:"lockAddress"`
	LockMethod           *string `pulumi:"lockMethod"`
	UnlockAddress        *string `pulumi:"unlockAddress"`
	UnlockMethod         *string `pulumi:"unlockMethod"`
	Username             *string `pulumi:"username"`
	Password             *string `pulumi:"password"`
	SkipCertVerification *bool   `pulumi:"skipCertVerification"`
	Workspace            *string `pulumi:"workspace"`
}

func (HttpStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*httpStateReferenceArgs)(nil)).Elem()
}
