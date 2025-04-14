// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Terraform.State.Inputs
{

    public sealed class Workspaces : global::Pulumi.InvokeArgs
    {
        /// <summary>
        /// The full name of one remote workspace. When configured, only the default workspace can be used. This option conflicts with prefix.
        /// </summary>
        [Input("name")]
        public string? Name { get; set; }

        /// <summary>
        /// A prefix used in the names of one or more remote workspaces, all of which can be used with this configuration. The full workspace names are used in HCP Terraform, and the short names (minus the prefix) are used on the command line for Terraform CLI workspaces. If omitted, only the default workspace can be used. This option conflicts with name.
        /// </summary>
        [Input("prefix")]
        public string? Prefix { get; set; }

        public Workspaces()
        {
        }
        public static new Workspaces Empty => new Workspaces();
    }
}
