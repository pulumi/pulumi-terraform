// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the etcd v2 backend. Note
    /// that there is a separate configuration class for state stored in the ectd v3 backend.
    /// </summary>
    public class EtcdV2RemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "etcd";
         
         /// <summary>
         /// The path at which to store the state.
         /// </summary>
         [Input("path", required: true)]
         public Input<string> Path { get; set; } = null!;
         
         /// <summary>
         /// A space-separated list of the etcd endpoints.
         /// </summary>
         [Input("endpoints", required: true)]
         public Input<string> Endpoints { get; set; } = null!;

         /// <summary>
         /// The username with which to authenticate to etcd.
         /// </summary>
         [Input("username")]
         public Input<string>? Username { get; set; }

         /// <summary>
         /// The password with which to authenticate to etcd.
         /// </summary>
         [Input("password")]
         public Input<string>? Password { get; set; }

         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
