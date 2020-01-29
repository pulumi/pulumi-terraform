// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the local enhanced
    /// backend.
    /// </summary>
    public class LocalBackendRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
        /// <summary>
        /// A constant describing the name of the Terraform backend, used as the discriminant
        /// for the union of backend configurations.
        /// </summary>
        [Input("backendType", required: true)]
        public override Input<string> BackendType => "local";

        /// <summary>
        /// The path to the Terraform state file.
        /// </summary>
        [Input("path", required: true)]
        public Input<string> Path { get; set; } = null!;
    }
}
