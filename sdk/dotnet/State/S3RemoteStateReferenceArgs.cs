// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the S3 backend.
    /// </summary>
    public class S3RemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
         /// <summary>
         /// A constant describing the name of the Terraform backend, used as the discriminant
         /// for the union of backend configurations.
         /// </summary>
         [Input("backendType", required: true)]
         public override Input<string> BackendType => "s3";
         
         /// <summary>
         /// The name of the S3 bucket.
         /// </summary>
         [Input("bucket", required: true)]
         public Input<string> Bucket { get; set; } = null!;
         
         /// <summary>
         /// The path to the state file inside the bucket. When using a non-default workspace,
         /// the state path will be `/workspace_key_prefix/workspace_name/key`.
         /// </summary>
         [Input("key", required: true)]
         public Input<string> Key { get; set; } = null!;
         
         /// <summary>
         /// The region of the S3 bucket. Also sourced from `AWS_DEFAULT_REGION` in the environment, if unset.
         /// </summary>
         [Input("region")]
         public Input<string>? Region { get; set; }

         /// <summary>
         /// A custom endpoint for the S3 API. Also sourced from `AWS_S3_ENDPOINT` in the environment, if unset.
         /// </summary>
         [Input("endpoint")]
         public Input<string>? Endpoint { get; set; }

         /// <summary>
         /// AWS Access Key. Sourced from the standard credentials pipeline, if unset.
         /// </summary>
         [Input("accessKey")]
         public Input<string>? AccessKey { get; set; }

         /// <summary>
         /// AWS Secret Access Key. Sourced from the standard credentials pipeline, if unset.
         /// </summary>
         [Input("secretKey")]
         public Input<string>? SecretKey { get; set; }

         /// <summary>
         /// The AWS profile name as set in the shared credentials file.
         /// </summary>
         [Input("profile")]
         public Input<string>? Profile { get; set; }

         /// <summary>
         /// The path to the shared credentials file. If this is not set and a profile is
         /// specified, `~/.aws/credentials` will be used by default.
         /// </summary>
         [Input("sharedCredentialsFile")]
         public Input<string>? SharedCredentialsFile { get; set; }

         /// <summary>
         /// An MFA token. Sourced from the `AWS_SESSION_TOKEN` in the environment variable if needed and unset.
         /// </summary>
         [Input("token")]
         public Input<string>? Token { get; set; }

         /// <summary>
         /// The ARN of an IAM Role to be assumed in order to read the state from S3.
         /// </summary>
         [Input("roleArn")]
         public Input<string>? RoleArn { get; set; }

         /// <summary>
         /// The external ID to use when assuming the IAM role.
         /// </summary>
         [Input("externalId")]
         public Input<string>? ExternalId { get; set; }

         /// <summary>
         /// The session name to use when assuming the IAM role.
         /// </summary>
         [Input("sessionName")]
         public Input<string>? SessionName { get; set; }

         /// <summary>
         /// The prefix applied to the state path inside the bucket. This is only relevant when
         /// using a non-default workspace, and defaults to `env:`.
         /// </summary>
         [Input("workspaceKeyPrefix")]
         public Input<string>? WorkspaceKeyPrefix { get; set; }

         /// <summary>
         /// A custom endpoint for the IAM API. Sourced from `AWS_IAM_ENDPOINT`, if unset.
         /// </summary>
         [Input("iamEndpoint")]
         public Input<string>? IamEndpoint { get; set; }

         /// <summary>
         /// A custom endpoint for the STS API. Sourced from `AWS_STS_ENDPOINT`, if unset.
         /// </summary>
         [Input("stsEndpoint")]
         public Input<string>? StsEndpoint { get; set; }

         /// <summary>
         /// The Terraform workspace from which to read state.
         /// </summary>
         [Input("workspace")]
         public Input<string>? Workspace { get; set; }
    }
}
