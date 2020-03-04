// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the remote enhanced
    /// backend.
    /// </summary>
    public class RemoteBackendRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "remote";
         
         /// <summary>
         /// The name of the organization containing the targeted workspace(s).
         /// </summary>
         [Input("organization", required: true)]
         public Input<string> Organization { get; set; } = null!;
         
         /// <summary>
         /// The remote backend hostname to which to connect. Defaults to `app.terraform.io`.
         /// </summary>
         [Input("hostname")]
         public Input<string>? Hostname { get; set; }

         /// <summary>
         /// The token used to authenticate with the remote backend.
         /// </summary>
         [Input("token")]
         public Input<string>? Token { get; set; }

         /// <summary>
         /// A block specifying which remote cdworkspace(s) to use.
         /// </summary>
         [Input("workspaces")]
         public Input<RemoteBackendWorkspaceConfig>? Workspaces { get; set; }
    }

    /// <summary>
    /// Configuration options for a workspace for use with the remote enhanced backend.
    /// </summary>
    public class RemoteBackendWorkspaceConfig : ResourceArgs
    {
        /// <summary>
        /// The full name of one remote workspace. When configured, only the default workspace
        /// can be used. This option conflicts with prefix.
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        /// <summary>
        /// A prefix used in the names of one or more remote workspaces, all of which can be used
        /// with this configuration. If unset, only the default workspace can be used. This option
        /// conflicts with name.
        /// </summary>
        [Input("prefix")]
        public Input<string>? Prefix { get; set; }
    }
}
