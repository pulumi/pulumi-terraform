// Copyright 2016-2020, Pulumi Corporation.

using System.Collections.Immutable;
using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// Manages a reference to a Terraform Remote State. The root outputs of the remote state are available
    /// via the <see cref="Outputs"/> property or the <see cref="GetOutput"/> method.
    /// </summary>
    public class RemoteStateReference : CustomResource
    {
        /// <summary>
        /// The root outputs of the referenced Terraform state. 
        /// </summary>
        [Output("outputs")]
        public Output<ImmutableDictionary<string, object>> Outputs { get; private set; } = null!;

        /// <summary>
        /// Create a RemoteStateReference resource with the given unique name, arguments, and options.
        /// </summary>
        /// <param name="name">The unique name of the remote state reference.</param>
        /// <param name="args">The arguments to use to populate this resource's properties.</param>
        /// <param name="options">A bag of options that control this resource's behavior.</param>
        public RemoteStateReference(string name, RemoteStateReferenceArgs args, CustomResourceOptions? options = null)
            : base("terraform:state:RemoteStateReference", 
                name, 
                args, 
                CustomResourceOptions.Merge(options, new CustomResourceOptions { Id = name }))
        {
        }

        /// <summary>
        /// Fetches the value of a root output from the Terraform Remote State.
        /// </summary>
        /// <param name="name">The name of the output to fetch. The name is formatted exactly as per
        /// the "output" block in the Terraform configuration.</param>
        /// <returns></returns>
        public Output<object> GetOutput(Input<string> name)
            => Output.Tuple(name.ToOutput(), Outputs).Apply(v => v.Item2[v.Item1]);
    }
    
    /// <summary>
    /// The base type for arguments for constructing a RemoteStateReference resource.
    /// </summary>
    public abstract class RemoteStateReferenceArgs : ResourceArgs
    {
        /// <summary>
        /// A constant describing the name of the Terraform backend, used as the discriminant
        /// for the union of backend configurations.
        /// </summary>
        public abstract Input<string> BackendType { get; }
    }
}
