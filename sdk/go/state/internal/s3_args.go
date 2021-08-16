package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type S3StateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	Bucket                    pulumi.StringInput
	Key                       pulumi.StringInput
	Region                    pulumi.StringPtrInput
	Endpoint                  pulumi.StringPtrInput
	AccessKey                 pulumi.StringPtrInput
	SecretKey                 pulumi.StringPtrInput
	Profile                   pulumi.StringPtrInput
	SharedCredentialsFile     pulumi.StringPtrInput
	Token                     pulumi.StringPtrInput
	RoleArn                   pulumi.StringPtrInput
	ExternalID                pulumi.StringPtrInput
	SessionName               pulumi.StringPtrInput
	WorkspaceKeyPrefix        pulumi.StringPtrInput
	IAMEndpoint               pulumi.StringPtrInput
	STSEndpoint               pulumi.StringPtrInput
	Workspace                 pulumi.StringPtrInput
	SkipRegionValidation      pulumi.BoolPtrInput
	SkipCredentialsValidation pulumi.BoolPtrInput
	SkipMetadataApiCheck      pulumi.BoolPtrInput
	ForcePathStyle            pulumi.BoolPtrInput
}

type s3StateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	Bucket                    string  `pulumi:"bucket"`
	Key                       string  `pulumi:"key"`
	Region                    *string `pulumi:"region"`
	Endpoint                  *string `pulumi:"endpoint"`
	AccessKey                 *string `pulumi:"accessKey"`
	SecretKey                 *string `pulumi:"secretKey"`
	Profile                   *string `pulumi:"profile"`
	SharedCredentialsFile     *string `pulumi:"sharedCredentialsFile"`
	Token                     *string `pulumi:"token"`
	RoleArn                   *string `pulumi:"roleArn"`
	ExternalID                *string `pulumi:"externalId"`
	SessionName               *string `pulumi:"sessionName"`
	WorkspaceKeyPrefix        *string `pulumi:"workspaceKeyPrefix"`
	IAMEndpoint               *string `pulumi:"iamEndpoint"`
	STSEndpoint               *string `pulumi:"stsEndpoint"`
	Workspace                 *string `pulumi:"workspace"`
	SkipRegionValidation      *bool   `pulumi:"skipRegionValidation"`
	SkipCredentialsValidation *bool   `pulumi:"skipCredentialsValidation"`
	SkipMetadataApiCheck      *bool   `pulumi:"skipMetadataApiCheck"`
	ForcePathStyle            *bool   `pulumi:"forcePathStyle"`
}

func (S3StateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*s3StateReferenceArgs)(nil)).Elem()
}
