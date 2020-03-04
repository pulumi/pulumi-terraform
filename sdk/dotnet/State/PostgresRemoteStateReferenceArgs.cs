// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the Postgres backend.
    /// </summary>
    public class PostgresRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "pg";
         
         /// <summary>
         /// Postgres connection string; a `postgres://` URL.
         /// </summary>
         [Input("connStr", required: true)]
         public Input<string> ConnStr { get; set; } = null!;
         
         /// <summary>
         /// Name of the automatically-managed Postgres schema. Defaults to `terraform_remote_state`.
         /// </summary>
         [Input("schemaName")]
         public Input<string>? SchemaName { get; set; }

         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
