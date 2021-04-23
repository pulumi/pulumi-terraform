package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type OssStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Bucket                pulumi.StringInput
	AccessKey             pulumi.StringPtrInput
	SecretKey             pulumi.StringPtrInput
	SecurityToken         pulumi.StringPtrInput
	EcsRoleName           pulumi.StringPtrInput
	Key                   pulumi.StringPtrInput
	Prefix                pulumi.StringPtrInput
	Region                pulumi.StringPtrInput
	Endpoint              pulumi.StringPtrInput
	Profile               pulumi.StringPtrInput
	SharedCredentialsFile pulumi.StringPtrInput
	RoleArn               pulumi.StringPtrInput
	Policy                pulumi.StringPtrInput
	SessionName           pulumi.StringPtrInput
	SessionExpiration     pulumi.StringPtrInput
}

type ossStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Bucket                string  `pulumi:"bucket"`
	AccessKey             *string `pulumi:"accessKey"`
	SecretKey             *string `pulumi:"secretKey"`
	SecurityToken         *string `pulumi:"securityToken"`
	EcsRoleName           *string `pulumi:"ecsRoleName"`
	Key                   *string `pulumi:"key"`
	Prefix                *string `pulumi:"prefix"`
	Region                *string `pulumi:"region"`
	Endpoint              *string `pulumi:"endpoint"`
	Profile               *string `pulumi:"profile"`
	SharedCredentialsFile *string `pulumi:"sharedCredentialsFile"`
	RoleArn               *string `pulumi:"roleArn"`
	Policy                *string `pulumi:"policy"`
	SessionName           *string `pulumi:"sessionName"`
	SessionExpiration     *string `pulumi:"sessionExpiration"`
}

func (OssStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*ossStateReferenceArgs)(nil)).Elem()
}
