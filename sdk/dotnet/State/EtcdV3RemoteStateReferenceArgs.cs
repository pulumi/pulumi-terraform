// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the etcd v3 backend. Note
    /// that there is a separate configuration class for state stored in the ectd v2 backend.
    /// </summary>
    public class EtcdV3RemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "etcdv3";
         
         [Input("endpoints", required: true)]
         private InputList<string>? _endpoints;

         /// <summary>
         /// A list of the etcd endpoints.
         /// </summary>
         public InputList<string> Endpoints
         {
             get => _endpoints ?? (_endpoints = new InputList<string>());
             set => _endpoints = value;
         }

         /// <summary>
         /// The username with which to authenticate to etcd. Sourced from `ETCDV3_USERNAME` in
         /// the environment, if unset.
         /// </summary>
         [Input("username")]
         public Input<string>? Username { get; set; }

         /// <summary>
         /// The password with which to authenticate to etcd. Sourced from `ETCDV3_PASSWORD` in
         /// the environment, if unset.
         /// </summary>
         [Input("password")]
         public Input<string>? Password { get; set; }

         /// <summary>
         /// An optional prefix to be added to keys when storing state in etcd.
         /// </summary>
         [Input("prefix")]
         public Input<string>? Prefix { get; set; }

         /// <summary>
         /// Path to a PEM-encoded certificate authority bundle with which to verify certificates
         /// of TLS-enabled etcd servers. 
         /// </summary>
         [Input("cacertPath")]
         public Input<string>? CACertPath { get; set; }

         /// <summary>
         /// Path to a PEM-encoded certificate to provide to etcd for client authentication.
         /// </summary>
         [Input("certPath")]
         public Input<string>? CertPath { get; set; }

         /// <summary>
         /// Path to a PEM-encoded key to use for client authentication.
         /// </summary>
         [Input("keyPath")]
         public Input<string>? KeyPath { get; set; }

         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
