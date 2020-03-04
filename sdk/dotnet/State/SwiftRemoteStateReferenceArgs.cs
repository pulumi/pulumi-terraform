// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the Swift backend.
    /// </summary>
    public class SwiftRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "swift";
         
         /// <summary>
         /// The Identity authentication URL. Sourced from `OS_AUTH_URL` in the environment, if unset.
         /// </summary>
         [Input("authUrl", required: true)]
         public Input<string> AuthUrl { get; set; } = null!;
         
         /// <summary>
         /// The name of the container in which the Terraform state file is stored.
         /// </summary>
         [Input("container", required: true)]
         public Input<string> Container { get; set; } = null!;
         
         /// <summary>
         /// The username with which to log in. Sourced from `OS_USERNAME` in the environment, if
         /// unset.
         /// </summary>
         [Input("userName")]
         public Input<string>? UserName { get; set; }

         /// <summary>
         /// The user ID with which to log in. Sourced from `OS_USER_ID` in the environment, if
         /// unset.  
         /// </summary>
         [Input("userId")]
         public Input<string>? UserId { get; set; }

         /// <summary>
         /// The password with which to log in. Sourced from `OS_PASSWORD` in the environment,
         /// if unset.
         /// </summary>
         [Input("password")]
         public Input<string>? Password { get; set; }

         /// <summary>
         /// Access token with which to log in in stead of a username and password. Sourced from
         /// `OS_AUTH_TOKEN` in the environment, if unset.
         /// </summary>
         [Input("token")]
         public Input<string>? Token { get; set; }

         /// <summary>
         /// The region in which the state file is stored. Sourced from `OS_REGION_NAME`, if
         /// unset.
         /// </summary>
         [Input("regionName", required: true)]
         public Input<string> RegionName { get; set; } = null!;

         /// <summary>
         /// The ID of the tenant (for identity v2) or project (identity v3) which which to log in.
         /// Sourced from `OS_TENANT_ID` or `OS_PROJECT_ID` in the environment, if unset.
         /// </summary>
         [Input("tenantId")]
         public Input<string>? TenantId { get; set; }

         /// <summary>
         /// The name of the tenant (for identity v2) or project (identity v3) which which to log in.
         /// Sourced from `OS_TENANT_NAME` or `OS_PROJECT_NAME` in the environment, if unset.
         /// </summary>
         [Input("tenantName")]
         public Input<string>? TenantName { get; set; }

         /// <summary>
         /// The ID of the domain to scope the log in to (identity v3). Sourced from `OS_USER_DOMAIN_ID`,
         /// `OS_PROJECT_DOMAIN_ID` or `OS_DOMAIN_ID` in the environment, if unset.
         /// </summary>
         [Input("domainId")]
         public Input<string>? DomainId { get; set; }

         /// <summary>
         /// The name of the domain to scope the log in to (identity v3). Sourced from
         /// `OS_USER_DOMAIN_NAME`, `OS_PROJECT_DOMAIN_NAME` or `OS_DOMAIN_NAME` in the environment,
         /// if unset. 
         /// </summary>
         [Input("domainName")]
         public Input<string>? DomainName { get; set; }

         /// <summary>
         /// Whether to disable verification of the server TLS certificate. Sourced from
         /// `OS_INSECURE` in the environment, if unset. 
         /// </summary>
         [Input("insecure")]
         public Input<bool>? Insecure { get; set; }

         /// <summary>
         /// A path to a CA root certificate for verifying the server TLS certificate. Sourced from
         /// `OS_CACERT` in the environment, if unset.
         /// </summary>
         [Input("cacertFile")]
         public Input<string>? CACertFile { get; set; }

         /// <summary>
         /// A path to a client certificate for TLS client authentication. Sourced from `OS_CERT`
         /// in the environment, if unset.
         /// </summary>
         [Input("cert")]
         public Input<string>? Cert { get; set; }

         /// <summary>
         /// A path to the private key corresponding to the client certificate for TLS client
         /// authentication. Sourced from `OS_KEY` in the environment, if unset.
         /// </summary>
         [Input("key")]
         public Input<string>? Key { get; set; }
    }
}
