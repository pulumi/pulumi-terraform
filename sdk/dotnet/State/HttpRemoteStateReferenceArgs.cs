// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    ///  The configuration options for a Terraform Remote State stored in the HTTP backend.
    /// </summary>
    public class HttpRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "http";
         
         /// <summary>
         /// The address of the HTTP endpoint.
         /// </summary>
         [Input("address", required: true)]
         public Input<string> address { get; set; } = null!;

         /// <summary>
         /// HTTP method to use when updating state. Defaults to `POST`.
         /// </summary>
         [Input("updateMethod")]
         public Input<string>? UpdateMethod { get; set; }

         /// <summary>
         /// The address of the lock REST endpoint. Not setting a value disables locking.
         /// </summary>
         [Input("lockAddress")]
         public Input<string>? LockAddress { get; set; }

         /// <summary>
         /// The HTTP method to use when locking. Defaults to `LOCK`.
         /// </summary>
         [Input("lockMethod")]
         public Input<string>? LockMethod { get; set; }

         /// <summary>
         /// The address of the unlock REST endpoint. Not setting a value disables locking.
         /// </summary>
         [Input("unlockAddress")]
         public Input<string>? UnlockAddress { get; set; }

         /// <summary>
         /// The HTTP method to use when unlocking. Defaults to `UNLOCK`.
         /// </summary>
         [Input("unlockMethod")]
         public Input<string>? UnlockMethod { get; set; }

         /// <summary>
         /// The username used for HTTP basic authentication.
         /// </summary>
         [Input("username")]
         public Input<string>? Username { get; set; }

         /// <summary>
         /// The password used for HTTP basic authentication.
         /// </summary>
         [Input("password")]
         public Input<string>? Password { get; set; }

         /// <summary>
         /// Whether to skip TLS verification. Defaults to false.
         /// </summary>
         [Input("skipCertVerification")]
         public Input<bool>? SkipCertVerification { get; set; }

         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
