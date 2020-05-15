package internal

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type AzureRMStateReferenceArgs struct {
	BackendType pulumi.StringPtrInput

	StorageAccountName pulumi.StringInput
	ContainerName      pulumi.StringInput
	Key                pulumi.StringPtrInput
	Environment        pulumi.StringPtrInput
	Endpoint           pulumi.StringPtrInput
	UseMSI             pulumi.BoolPtrInput
	SubscriptionID     pulumi.StringPtrInput
	TenantID           pulumi.StringPtrInput
	MSIEndpoint        pulumi.StringPtrInput
	SasToken           pulumi.StringPtrInput
	AccessKey          pulumi.StringPtrInput
	ResourceGroupName  pulumi.StringPtrInput
	ClientID           pulumi.StringPtrInput
	ClientSecret       pulumi.StringPtrInput
	Workspace          pulumi.StringPtrInput
}

type azureRMStateReferenceArgs struct {
	BackendType *string `pulumi:"backendType"`

	StorageAccountName string  `pulumi:"storageAccountName"`
	ContainerName      string  `pulumi:"containerName"`
	Key                *string `pulumi:"key"`
	Environment        *string `pulumi:"environment"`
	Endpoint           *string `pulumi:"endpoint"`
	UseMSI             *bool   `pulumi:"useMsi"`
	SubscriptionID     *string `pulumi:"subscriptionId"`
	TenantID           *string `pulumi:"tenantId"`
	MSIEndpoint        *string `pulumi:"msiEndpoint"`
	SasToken           *string `pulumi:"sasToken"`
	AccessKey          *string `pulumi:"accessKey"`
	ResourceGroupName  *string `pulumi:"resourceGroupName"`
	ClientID           *string `pulumi:"clientId"`
	ClientSecret       *string `pulumi:"clientSecret"`
	Workspace          *string `pulumi:"workspace"`
}

func (AzureRMStateReferenceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*azureRMStateReferenceArgs)(nil)).Elem()
}
