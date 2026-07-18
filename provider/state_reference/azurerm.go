// Copyright 2016-2025, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"

	"github.com/hashicorp/terraform/shim"
	"github.com/zclconf/go-cty/cty"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type GetAzureRMReference struct{}

var (
	_ = (infer.Annotated)((*GetAzureRMReference)(nil))
	_ = (infer.ExplicitDependencies[GetAzureRMReferenceArgs, StateReferenceOutputs])((*GetAzureRMReference)(nil))
)

func (r *GetAzureRMReference) Annotate(a infer.Annotator) {
	a.Describe(&r, "Access state stored in an Azure Blob Storage container.")
}

// Taken from https://developer.hashicorp.com/terraform/language/backend/azurerm#configuration-variables
//
// Only the arguments that affect reading state are exposed. Write-only arguments
// (such as snapshot) have no effect on a read and are omitted.
type GetAzureRMReferenceArgs struct {
	StorageAccountName string  `pulumi:"storageAccountName"`
	ContainerName      string  `pulumi:"containerName"`
	Key                string  `pulumi:"key"`
	Workspace          *string `pulumi:"workspace,optional"`

	Environment  *string `pulumi:"environment,optional"`
	MetadataHost *string `pulumi:"metadataHost,optional"`
	Endpoint     *string `pulumi:"endpoint,optional"`

	AccessKey         *string `pulumi:"accessKey,optional" provider:"secret"`
	SasToken          *string `pulumi:"sasToken,optional" provider:"secret"`
	ResourceGroupName *string `pulumi:"resourceGroupName,optional"`

	SubscriptionID *string `pulumi:"subscriptionId,optional"`
	TenantID       *string `pulumi:"tenantId,optional"`
	ClientID       *string `pulumi:"clientId,optional"`

	ClientSecret              *string `pulumi:"clientSecret,optional" provider:"secret"`
	ClientCertificatePath     *string `pulumi:"clientCertificatePath,optional"`
	ClientCertificatePassword *string `pulumi:"clientCertificatePassword,optional" provider:"secret"`

	UseMsi      *bool   `pulumi:"useMsi,optional"`
	MsiEndpoint *string `pulumi:"msiEndpoint,optional"`

	UseOidc           *bool   `pulumi:"useOidc,optional"`
	OidcToken         *string `pulumi:"oidcToken,optional" provider:"secret"`
	OidcTokenFilePath *string `pulumi:"oidcTokenFilePath,optional"`
	OidcRequestURL    *string `pulumi:"oidcRequestUrl,optional"`
	OidcRequestToken  *string `pulumi:"oidcRequestToken,optional" provider:"secret"`

	UseAzureadAuth *bool `pulumi:"useAzureadAuth,optional"`
}

func (r *GetAzureRMReferenceArgs) Annotate(a infer.Annotator) {
	a.Describe(&r.StorageAccountName, "The name of the storage account.")
	a.Describe(&r.ContainerName, "The name of the storage container within the storage account.")
	a.Describe(&r.Key, "The name of the blob holding the Terraform state file inside the storage container.")
	a.Describe(&r.Workspace, "The Terraform workspace to read state from.")

	a.Describe(&r.Environment, "The Azure cloud environment to use: public (default), china, "+
		"german, stack or usgovernment. Falls back to the ARM_ENVIRONMENT environment variable when unset.")
	a.Describe(&r.MetadataHost, "The hostname of the Azure metadata service used to obtain the cloud "+
		"environment. Falls back to the ARM_METADATA_HOST environment variable when unset.")
	a.Describe(&r.Endpoint, "A custom endpoint for the Azure Resource Manager API. Falls back to the "+
		"ARM_ENDPOINT environment variable when unset.")

	a.Describe(&r.AccessKey, "The access key of the storage account. Falls back to the ARM_ACCESS_KEY "+
		"environment variable when unset.")
	a.Describe(&r.SasToken, "A SAS token for accessing the storage container. Falls back to the "+
		"ARM_SAS_TOKEN environment variable when unset.")
	a.Describe(&r.ResourceGroupName, "The name of the resource group holding the storage account. Required "+
		"when using AzureAD authentication against the Azure Resource Manager API to look up the storage access key.")

	a.Describe(&r.SubscriptionID, "The subscription ID holding the storage account. Falls back to the "+
		"ARM_SUBSCRIPTION_ID environment variable when unset.")
	a.Describe(&r.TenantID, "The tenant ID to authenticate against. Falls back to the ARM_TENANT_ID "+
		"environment variable when unset.")
	a.Describe(&r.ClientID, "The client ID to authenticate as. Falls back to the ARM_CLIENT_ID "+
		"environment variable when unset.")

	a.Describe(&r.ClientSecret, "The client secret used for service principal authentication. Falls back "+
		"to the ARM_CLIENT_SECRET environment variable when unset.")
	a.Describe(&r.ClientCertificatePath, "The path to the PFX file used as the client certificate for "+
		"service principal authentication. Falls back to the ARM_CLIENT_CERTIFICATE_PATH environment "+
		"variable when unset.")
	a.Describe(&r.ClientCertificatePassword, "The password for the client certificate specified in "+
		"clientCertificatePath. Falls back to the ARM_CLIENT_CERTIFICATE_PASSWORD environment variable when unset.")

	a.Describe(&r.UseMsi, "Whether to authenticate using Managed Service Identity. Falls back to the "+
		"ARM_USE_MSI environment variable when unset.")
	a.Describe(&r.MsiEndpoint, "The endpoint of the Managed Service Identity. Falls back to the "+
		"ARM_MSI_ENDPOINT environment variable when unset.")

	a.Describe(&r.UseOidc, "Whether to authenticate using OIDC. Falls back to the ARM_USE_OIDC "+
		"environment variable when unset.")
	a.Describe(&r.OidcToken, "A JWT token for OIDC authentication. Conflicts with oidcRequestToken. Falls "+
		"back to the ARM_OIDC_TOKEN environment variable when unset.")
	a.Describe(&r.OidcTokenFilePath, "The path to a file containing a JWT token for OIDC authentication. "+
		"Conflicts with oidcRequestToken. Falls back to the ARM_OIDC_TOKEN_FILE_PATH environment variable when unset.")
	a.Describe(&r.OidcRequestURL, "The URL of the OIDC provider to request an ID token from, e.g. in "+
		"GitHub Actions. Requires oidcRequestToken. Falls back to the ARM_OIDC_REQUEST_URL or "+
		"ACTIONS_ID_TOKEN_REQUEST_URL environment variables when unset.")
	a.Describe(&r.OidcRequestToken, "The bearer token for requests to the oidcRequestUrl URL. Falls back "+
		"to the ARM_OIDC_REQUEST_TOKEN or ACTIONS_ID_TOKEN_REQUEST_TOKEN environment variables when unset.")

	a.Describe(&r.UseAzureadAuth, "Whether to authenticate against the storage container with AzureAD "+
		"instead of an access key. Falls back to the ARM_USE_AZUREAD environment variable when unset.")

	a.SetDefault(&r.Workspace, defaultWorkspace)
}

// WireDependencies lets us tell users that our outputs shouldn't be secret, even when
// the credentials (when provided) are always secret.
//
// TODO[https://github.com/pulumi/pulumi-go-provider/issues/323]: This doesn't currently
// work; [infer.ExplicitDependencies] is not currently implemented for [infer] based
// functions.
func (r *GetAzureRMReference) WireDependencies(
	f infer.FieldSelector, _ *GetAzureRMReferenceArgs, state *StateReferenceOutputs,
) {
	f.OutputField(&state).NeverSecret() // The output should never be secret by default
}

// backendConfig builds the azurerm backend configuration, keyed by the backend's
// attribute names.
func (r *GetAzureRMReferenceArgs) backendConfig() map[string]cty.Value {
	return map[string]cty.Value{
		"storage_account_name":        cty.StringVal(r.StorageAccountName),
		"container_name":              cty.StringVal(r.ContainerName),
		"key":                         cty.StringVal(r.Key),
		"environment":                 ctyStringOrNil(r.Environment),
		"metadata_host":               ctyStringOrNil(r.MetadataHost),
		"endpoint":                    ctyStringOrNil(r.Endpoint),
		"access_key":                  ctyStringOrNil(r.AccessKey),
		"sas_token":                   ctyStringOrNil(r.SasToken),
		"resource_group_name":         ctyStringOrNil(r.ResourceGroupName),
		"subscription_id":             ctyStringOrNil(r.SubscriptionID),
		"tenant_id":                   ctyStringOrNil(r.TenantID),
		"client_id":                   ctyStringOrNil(r.ClientID),
		"client_secret":               ctyStringOrNil(r.ClientSecret),
		"client_certificate_path":     ctyStringOrNil(r.ClientCertificatePath),
		"client_certificate_password": ctyStringOrNil(r.ClientCertificatePassword),
		"use_msi":                     ctyBoolOrNil(r.UseMsi),
		"msi_endpoint":                ctyStringOrNil(r.MsiEndpoint),
		"use_oidc":                    ctyBoolOrNil(r.UseOidc),
		"oidc_token":                  ctyStringOrNil(r.OidcToken),
		"oidc_token_file_path":        ctyStringOrNil(r.OidcTokenFilePath),
		"oidc_request_url":            ctyStringOrNil(r.OidcRequestURL),
		"oidc_request_token":          ctyStringOrNil(r.OidcRequestToken),
		"use_azuread_auth":            ctyBoolOrNil(r.UseAzureadAuth),
	}
}

func (r *GetAzureRMReference) Invoke(
	ctx context.Context, req infer.FunctionRequest[GetAzureRMReferenceArgs],
) (infer.FunctionResponse[StateReferenceOutputs], error) {
	args := req.Input

	results, err := shim.StateReferenceRead(ctx, "azurerm", *args.Workspace, args.backendConfig())

	return infer.FunctionResponse[StateReferenceOutputs]{Output: StateReferenceOutputs{results}}, err
}
