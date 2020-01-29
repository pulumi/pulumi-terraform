// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the AzureRM backend. 
    /// </summary>
    public class AzureRMRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "azurerm";
         
         /// <summary>
         /// The name of the storage account.
         /// </summary>
         [Input("storageAccountName", required: true)]
         public Input<string> StorageAccountName { get; set; } = null!;

         /// <summary>
         /// The name of the storage container within the storage account.
         /// </summary>
         [Input("containerName", required: true)]
         public Input<string> ContainerName { get; set; } = null!;

         /// <summary>
         /// The name of the blob in representing the Terraform State file inside the storage container.
         /// </summary>
         [Input("key")]
         public Input<string>? Key { get; set; }
         
         /// <summary>
         /// The Azure environment which should be used. Possible values are `public` (default), `china`,
         /// `german`, `stack` and `usgovernment`. Sourced from `ARM_ENVIRONMENT`, if unset.
         /// </summary>
         [Input("environment")]
         public Input<string>? Environment { get; set; }

         /// <summary>
         /// The custom endpoint for Azure Resource Manager. Sourced from `ARM_ENDPOINT`, if unset.
         /// </summary>
         [Input("endpoint")]
         public Input<string>? Endpoint { get; set; }
         
         /// <summary>
         /// Whether to authenticate using Managed Service Identity (MSI). Sourced from `ARM_USE_MSI`
         /// if unset. Defaults to false if no value is specified.
         /// </summary>
         [Input("useMsi")]
         public Input<bool>? UseMsi { get; set; }
         
         /// <summary>
         /// The Subscription ID in which the Storage Account exists. Used when authenticating using
         /// the Managed Service Identity (MSI) or a service principal. Sourced from `ARM_SUBSCRIPTION_ID`,
         /// if unset.
         /// </summary>
         [Input("subscriptionId")]
         public Input<string>? SubscriptionId { get; set; }

         /// <summary>
         /// The Tenant ID in which the Subscription exists. Used when authenticating using the
         /// Managed Service Identity (MSI) or a service principal. Sourced from `ARM_TENANT_ID`,
         /// if unset.
         /// </summary>
         [Input("tenantId")]
         public Input<string>? TenantId { get; set; }
         
         /// <summary>
         /// The path to a custom Managed Service Identity endpoint. Used when authenticating using
         /// the Managed Service Identity (MSI). Sourced from `ARM_MSI_ENDPOINT` in the environment,
         /// if unset. Automatically determined, if no value is provided.
         /// </summary>
         [Input("msiEndpoint")]
         public Input<string>? MsiEndpoint { get; set; }
         
         /// <summary>
         /// The SAS Token used to access the Blob Storage Account. Used when authenticating using
         /// a SAS Token. Sourced from `ARM_SAS_TOKEN` in the environment, if unset.
         /// </summary>
         [Input("sasToken")]
         public Input<string>? SasToken { get; set; }

         /// <summary>
         /// The Access Key used to access the blob storage account. Used when authenticating using
         /// an access key. Sourced from `ARM_ACCESS_KEY` in the environment, if unset.
         /// </summary>
         [Input("accessKey")]
         public Input<string>? AccessKey { get; set; }
         
         /// <summary>
         /// The name of the resource group in which the storage account exists. Used when authenticating
         /// using a service principal.
         /// </summary>
         [Input("resourceGroupName")]
         public Input<string>? ResourceGroupName { get; set; }
         
         /// <summary>
         /// The client ID of the service principal. Used when authenticating using a service principal.
         /// Sourced from `ARM_CLIENT_ID` in the environment, if unset.
         /// </summary>
         [Input("clientId")]
         public Input<string>? ClientId { get; set; }

         /// <summary>
         /// The client secret of the service principal. Used when authenticating using a service principal.
         /// Sourced from `ARM_CLIENT_SECRET` in the environment, if unset.
         /// </summary>
         [Input("clientSecret")]
         public Input<string>? ClientSecret { get; set; }

         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
