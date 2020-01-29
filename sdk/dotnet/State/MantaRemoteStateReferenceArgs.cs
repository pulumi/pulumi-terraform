// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the Manta backend.
    /// </summary>
    public class MantaRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "manta";
         
         /// <summary>
         /// The name of the Manta account. Sourced from `SDC_ACCOUNT` or `_ACCOUNT` in the
         /// environment, if unset.
         /// </summary>
         [Input("account", required: true)]
         public Input<string> Account { get; set; } = null!;
         
         /// <summary>
         /// The username of the Manta account with which to authenticate.
         /// </summary>
         [Input("user")]
         public Input<string>? User { get; set; }

         /// <summary>
         /// The Manta API Endpoint. Sourced from `MANTA_URL` in the environment, if unset.
         /// Defaults to `https://us-east.manta.joyent.com`.
         /// </summary>
         [Input("url")]
         public Input<string>? Url { get; set; }

         /// <summary>
         /// The private key material corresponding with the SSH key whose fingerprint is
         /// specified in keyId. Sourced from `SDC_KEY_MATERIAL` or `TRITON_KEY_MATERIAL`
         /// in the environment, if unset. If no value is specified, the local SSH agent
         /// is used for signing requests.
         /// </summary>
         [Input("keyMaterial")]
         public Input<string>? KeyMaterial { get; set; }
         
         /// <summary>
         /// The fingerprint of the public key matching the key material specified in
         /// keyMaterial, or in the local SSH agent.
         /// </summary>
         [Input("keyId", required: true)]
         public Input<string> KeyId { get; set; } = null!;
         
         /// <summary>
         /// The path relative to your private storage directory (`/$MANTA_USER/stor`)
         /// where the state file will be stored.
         /// </summary>
         [Input("path", required: true)]
         public Input<string> Path { get; set; } = null!;
         
         /// <summary>
         /// Skip verifying the TLS certificate presented by the Manta endpoint. This can
         /// be useful for installations which do not have a certificate signed by a trusted
         /// root CA. Defaults to false.
         /// </summary>
         [Input("insecureSkipTlsVerify", required: true)]
         public Input<bool> InsecureSkipTlsVerify { get; set; } = null!;
         
         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
