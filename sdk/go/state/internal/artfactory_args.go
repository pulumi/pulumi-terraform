package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ArtifactoryStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Username  pulumi.StringPtrInput
	Password  pulumi.StringPtrInput
	Url       pulumi.StringPtrInput
	Repo      pulumi.StringInput
	Subpath   pulumi.StringInput
	Workspace pulumi.StringPtrInput
}

type artifactoryStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Username  *string `pulumi:"username"`
	Password  *string `pulumi:"password"`
	Url       *string `pulumi:"url"`
	Repo      string  `pulumi:"repo"`
	Subpath   string  `pulumi:"subpath"`
	Workspace *string `pulumi:"workspace"`
}

func (ArtifactoryStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*artifactoryStateReferenceArgs)(nil)).Elem()
}
