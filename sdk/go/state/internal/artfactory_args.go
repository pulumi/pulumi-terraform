package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type ArtifatoryStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Username  pulumi.StringPtrInput
	Password  pulumi.StringPtrInput
	Url       pulumi.StringPtrInput
	Repo      pulumi.StringInput
	Subpath   pulumi.StringInput
	Workspace pulumi.StringPtrInput
}

type artifatoryStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Username  *string `pulumi:"username"`
	Password  *string `pulumi:"password"`
	Url       *string `pulumi:"url"`
	Repo      string  `pulumi:"repo"`
	Subpath   string  `pulumi:"subpath"`
	Workspace *string `pulumi:"workspace"`
}

func (ArtifatoryStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*artifatoryStateReferenceArgs)(nil)).Elem()
}
