package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type ConsulStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Path        pulumi.StringInput
	AccessToken pulumi.StringInput
	Address     pulumi.StringPtrInput
	Scheme      pulumi.StringPtrInput
	Datacenter  pulumi.StringPtrInput
	HttpAuth    pulumi.StringPtrInput
	Gzip        pulumi.BoolPtrInput
	CAFile      pulumi.StringPtrInput
	CertFile    pulumi.StringPtrInput
	KeyFile     pulumi.StringPtrInput
	Workspace   pulumi.StringPtrInput
}

type consulStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Path        string  `pulumi:"path"`
	AccessToken string  `pulumi:"accessToken"`
	Address     *string `pulumi:"address"`
	Scheme      *string `pulumi:"schema"`
	Datacenter  *string `pulumi:"datacenter"`
	HttpAuth    *string `pulumi:"httpAuth"`
	Gzip        *string `pulumi:"gzip"`
	CAFile      *string `pulumi:"caFile"`
	CertFile    *string `pulumi:"certFile"`
	KeyFile     *string `pulumi:"keyFile"`
	Workspace   *string `pulumi:"workspace"`
}

func (ConsulStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*consulStateReferenceArgs)(nil)).Elem()
}
