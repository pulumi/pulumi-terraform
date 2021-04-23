package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SwiftStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	AuthUrl    pulumi.StringInput
	Container  pulumi.StringInput
	UserName   pulumi.StringPtrInput
	UserID     pulumi.StringPtrInput
	Password   pulumi.StringPtrInput
	Token      pulumi.StringPtrInput
	RegionName pulumi.StringInput
	TenantID   pulumi.StringPtrInput
	TenantName pulumi.StringPtrInput
	DomainID   pulumi.StringPtrInput
	DomainName pulumi.StringPtrInput
	Insecure   pulumi.BoolPtrInput
	CACertFile pulumi.StringPtrInput
	Cert       pulumi.StringPtrInput
	Key        pulumi.StringPtrInput
}

type swiftStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	AuthUrl    string  `pulumi:"authUrl"`
	Container  string  `pulumi:"container"`
	UserName   *string `pulumi:"userName"`
	UserID     *string `pulumi:"userId"`
	Password   *string `pulumi:"password"`
	Token      *string `pulumi:"token"`
	RegionName string  `pulumi:"regionName"`
	TenantID   *string `pulumi:"tenantId"`
	TenantName *string `pulumi:"tenantName"`
	DomainID   *string `pulumi:"domainId"`
	DomainName *string `pulumi:"domainName"`
	Insecure   *bool   `pulumi:"insecure"`
	CACertFile *string `pulumi:"cacertFile"`
	Cert       *string `pulumi:"cert"`
	Key        *string `pulumi:"key"`
}

func (SwiftStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*swiftStateReferenceArgs)(nil)).Elem()
}
