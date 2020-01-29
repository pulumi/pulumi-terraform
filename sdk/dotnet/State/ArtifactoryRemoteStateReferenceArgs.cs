// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the Artifactory backend. 
    /// </summary>
    public class ArtifactoryRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "artifactory";
         
         /// <summary>
         /// The username with which to authenticate to Artifactory. Sourced from `ARTIFACTORY_USERNAME`
         /// in the environment, if unset.
         /// </summary>
         [Input("username")]
         public Input<string>? Username { get; set; }
         
         /// <summary>
         /// The password with which to authenticate to Artifactory. Sourced from `ARTIFACTORY_PASSWORD`
         /// in the environment, if unset.
         /// </summary>
         [Input("password")]
         public Input<string>? Password { get; set; }
         
         /// <summary>
         /// The Artifactory URL. Note that this is the base URL to artifactory, not the full repo and
         /// subpath. However, it must include the path to the artifactory installation - likely this
         /// will end in `/artifactory`. Sourced from `ARTIFACTORY_URL` in the environment, if unset.
         /// </summary>
         [Input("url")]
         public Input<string>? Url { get; set; }

         /// <summary>
         /// The repository name.
         /// </summary>
         [Input("repo", required: true)]
         public Input<string> Repo { get; set; } = null!;

         /// <summary>
         /// Path within the repository.
         /// </summary>
         [Input("subpath", required: true)]
         public Input<string> Subpath { get; set; } = null!;
         
         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
