// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the Consul backend. 
    /// </summary>
    public class ConsulRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "consul";
         
         /// <summary>
         /// Path in the Consul KV store.
         /// </summary>
         [Input("path", required: true)]
         public Input<string> Path { get; set; } = null!;

         /// <summary>
         /// Consul Access Token. Sourced from `CONSUL_HTTP_TOKEN` in the environment, if unset.
         /// </summary>
         [Input("accessToken", required: true)]
         public Input<string> AccessToken { get; set; } = null!;

         /// <summary>
         /// DNS name and port of the Consul HTTP endpoint specified in the format `dnsname:port`. Defaults
         /// to the local agent HTTP listener.
         /// </summary>
         [Input("address")]
         public Input<string>? Address { get; set; }

         /// <summary>
         /// Specifies which protocol to use when talking to the given address - either `http` or `https`. TLS
         /// support can also be enabled by setting the environment variable `CONSUL_HTTP_SSL` to `true`.
         /// </summary>
         [Input("scheme")]
         public Input<string>? Scheme { get; set; }

         /// <summary>
         /// The datacenter to use. Defaults to that of the agent.
         /// </summary>
         [Input("datacenter")]
         public Input<string>? Datacenter { get; set; }

         /// <summary>
         /// HTTP Basic Authentication credentials to be used when communicating with Consul, in the format of
         /// either `user` or `user:pass`. Sourced from `CONSUL_HTTP_AUTH`, if unset.
         /// </summary>
         [Input("httpAuth")]
         public Input<string>? HttpAuth { get; set; }
         
         /// <summary>
         /// Whether to compress the state data using gzip. Set to `true` to compress the state data using gzip,
         /// or `false` (default) to leave it uncompressed.
         /// </summary>
         [Input("gzip")]
         public Input<bool>? Gzip { get; set; }
         
         /// <summary>
         /// A path to a PEM-encoded certificate authority used to verify the remote agent's certificate. Sourced
         /// from `CONSUL_CAFILE` in the environment, if unset.
         /// </summary>
         [Input("caFile")]
         public Input<string>? CAFile { get; set; }
         
         /// <summary>
         /// A path to a PEM-encoded certificate provided to the remote agent; requires use of key_file. Sourced
         /// from `CONSUL_CLIENT_CERT` in the environment, if unset.
         /// </summary>
         [Input("certFile")]
         public Input<string>? CertFile { get; set; }
         
         /// <summary>
         /// A path to a PEM-encoded private key, required if cert_file is specified. Sourced from `CONSUL_CLIENT_KEY`
         /// in the environment, if unset.
         /// </summary>
         [Input("keyFile")]
         public Input<string>? KeyFile { get; set; }
         
         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
