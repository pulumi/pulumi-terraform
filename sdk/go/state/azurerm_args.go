package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v2/go/state/internal"
)

// AzureRMArgs specifies the configuration options for a Terraform Remote State
// stored in the AzureRM backend
type AzureRMArgs struct {
	// StorageAccountName is the name of the storage account
	StorageAccountName pulumi.StringInput

	// ContainerName is the name of the container within the storage account
	ContainerName pulumi.StringInput

	// Key is the name of the blob in representing the Terraform State file inside the storage container.
	Key pulumi.StringPtrInput

	// Environment is the Azure environment. Possible values are `public` (default), `china`,
	// `german`, `stack` and `usgovernment`. Sourced from `ARM_ENVIRONMENT`, if unset.
	Environment pulumi.StringPtrInput

	// Endpoint is custom endpoint for Azure Resource Manager. Sourced from `ARM_ENDPOINT`, if unset.
	Endpoint pulumi.StringPtrInput

	// UseMSI is whether to authenticate using Managed Service Identity (MSI). Sourced from `ARM_USE_MSI`
	// if unset. Defaults to `false` if no value is specified.
	UseMSI pulumi.BoolPtrInput

	// SubscriptionID is the subscription ID in which the Storage Account exists. Used when authenticating using
	// the Managed Service Identity (MSI) or a service principal. Sourced from `ARM_SUBSCRIPTION_ID`, if unset
	SubscriptionID pulumi.StringPtrInput

	// TenantID is the Tenant ID in which the Subscription exists. Used when authenticating using the
	// Managed Service Identity (MSI) or a service principal. Sourced from `ARM_TENANT_ID`, if unset.
	TenantID pulumi.StringPtrInput

	// MSIEndpoint is the path to a custom Managed Service Identity endpoint. Used when authenticating using
	// the Managed Service Identity (MSI). Sourced from `ARM_MSI_ENDPOINT` in the environment, if unset.
	// Automatically determined, if no value is provided.
	MSIEndpoint pulumi.StringPtrInput

	// SasToken is the SAS Token used to access the Blob Storage Account. Used when authenticating using
	// a SAS Token. Sourced from `ARM_SAS_TOKEN` in the environment, if unset.
	SasToken pulumi.StringPtrInput

	// AccessKey is the Access Key used to access the blob storage account. Used when authenticating using
	// an access key. Sourced from `ARM_ACCESS_KEY` in the environment, if unset.
	AccessKey pulumi.StringPtrInput

	// ResourceGroupName is the name of the resource group in which the storage account exists. Used when
	// authenticating using a service principal.
	ResourceGroupName pulumi.StringPtrInput

	// ClientID is client ID of the service principal. Used when authenticating using a service principal.
	// Sourced from `ARM_CLIENT_ID` in the environment, if unset.
	ClientID pulumi.StringPtrInput

	// ClientSecret is client secret of the service principal. Used when authenticating using a service principal.
	// Sourced from `ARM_CLIENT_SECRET` in the environment, if unset.
	ClientSecret pulumi.StringPtrInput

	// Workspace os the Terraform workspace from which to read state.
	Workspace pulumi.StringPtrInput
}

func (a *AzureRMArgs) toInternalArgs() pulumi.Input {
	return internal.AzureRMStateReferenceArgs{
		BackendType:        pulumi.String("azurerm"),
		StorageAccountName: a.StorageAccountName,
		ContainerName:      a.ContainerName,
		Key:                a.Key,
		Environment:        a.Environment,
		Endpoint:           a.Endpoint,
		UseMSI:             a.UseMSI,
		SubscriptionID:     a.SubscriptionID,
		TenantID:           a.TenantID,
		MSIEndpoint:        a.MSIEndpoint,
		SasToken:           a.SasToken,
		AccessKey:          a.AccessKey,
		ResourceGroupName:  a.ResourceGroupName,
		ClientID:           a.ClientID,
		ClientSecret:       a.ClientSecret,
		Workspace:          a.Workspace,
	}
}

func (l *AzureRMArgs) validateArgs() error {
	if l.StorageAccountName == pulumi.String("") || l.ContainerName == pulumi.String("") {
		return errors.New("`StorageAccountName` and `ContainerName` are required parameters")
	}
	return nil
}
