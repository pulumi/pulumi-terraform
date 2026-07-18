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
	"testing"

	"github.com/hashicorp/terraform/shim"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

// TestAzureRMBackendConfigMatchesSchema validates backendConfig against the real
// azurerm backend schema: every key must be a schema attribute of the right type,
// and all required attributes must be present. A key typo or type mismatch fails
// CoerceValue; a validation violation fails PrepareConfig.
func TestAzureRMBackendConfigMatchesSchema(t *testing.T) {
	InitTfBackend()
	backend := shim.BackendFactory("azurerm")()

	tests := []struct {
		name string
		args GetAzureRMReferenceArgs
	}{
		{
			name: "required only",
			args: GetAzureRMReferenceArgs{
				StorageAccountName: "testaccount",
				ContainerName:      "tfstate",
				Key:                "prod.terraform.tfstate",
			},
		},
		{
			name: "all arguments",
			args: GetAzureRMReferenceArgs{
				StorageAccountName:        "testaccount",
				ContainerName:             "tfstate",
				Key:                       "prod.terraform.tfstate",
				Environment:               ptr("public"),
				MetadataHost:              ptr("management.azure.com"),
				Endpoint:                  ptr("https://management.azure.com"),
				AccessKey:                 ptr("access-key"),
				SasToken:                  ptr("sas-token"),
				ResourceGroupName:         ptr("my-rg"),
				SubscriptionID:            ptr("00000000-0000-0000-0000-000000000000"),
				TenantID:                  ptr("00000000-0000-0000-0000-000000000000"),
				ClientID:                  ptr("00000000-0000-0000-0000-000000000000"),
				ClientSecret:              ptr("client-secret"),
				ClientCertificatePath:     ptr("/path/to/cert.pfx"),
				ClientCertificatePassword: ptr("cert-password"),
				UseMsi:                    ptr(true),
				MsiEndpoint:               ptr("http://169.254.169.254/metadata/identity/oauth2/token"),
				UseOidc:                   ptr(true),
				OidcToken:                 ptr("oidc-token"),
				OidcTokenFilePath:         ptr("/path/to/token"),
				OidcRequestURL:            ptr("https://token.actions.githubusercontent.com"),
				OidcRequestToken:          ptr("request-token"),
				UseAzureadAuth:            ptr(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coerced, err := backend.ConfigSchema().CoerceValue(cty.ObjectVal(tt.args.backendConfig()))
			require.NoError(t, err)

			_, diags := backend.PrepareConfig(coerced)
			require.False(t, diags.HasErrors(), "PrepareConfig: %v", diags.Err())
		})
	}
}
