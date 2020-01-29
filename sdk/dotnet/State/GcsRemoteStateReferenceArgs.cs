// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the Google Cloud Storage
    /// backend.
    /// </summary>
    public class GcsRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "gcs";
         
         /// <summary>
         /// The name of the Google Cloud Storage bucket.
         /// </summary>
         [Input("bucket", required: true)]
         public Input<string> Bucket { get; set; } = null!;
         
         /// <summary>
         /// Local path to Google Cloud Platform account credentials in JSON format. Sourced from
         /// `GOOGLE_CREDENTIALS` in the environment if unset. If no value is provided Google
         /// Application Default Credentials are used.
         /// </summary>
         [Input("credentials")]
         public Input<string>? Credentials { get; set; }
         
         /// <summary>
         /// Prefix used inside the Google Cloud Storage bucket. Named states for workspaces
         /// are stored in an object named `&lt;prefix&gt;/&lt;name&gt;.tfstate`.
         /// </summary>
         [Input("prefix")]
         public Input<string>? Prefix { get; set; }

         /// <summary>
         /// A 32 byte, base64-encoded customer supplied encryption key used to encrypt the
         /// state. Sourced from `GOOGLE_ENCRYPTION_KEY` in the environment, if unset.
         /// </summary>
         [Input("encryptionKey")]
         public Input<string>? EncryptionKey { get; set; }

         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
