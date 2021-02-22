// Copyright 2016-2020, Pulumi Corporation.

using Pulumi.Serialization;

namespace Pulumi.Terraform.State
{
    /// <summary>
    /// The configuration options for a Terraform Remote State stored in the OSS backend.
    /// </summary>
    public class OssRemoteStateReferenceArgs : RemoteStateReferenceArgs
    {
        /// <summary>
        /// A constant describing the name of the Terraform backend, used as the discriminant
        /// for the union of backend configurations.
        /// </summary>
        [Input("backendType", required: true)]
        public override Input<string> BackendType => "oss";

        /// <summary>
        /// The name of the Oss bucket.
        /// </summary>
        [Input("bucket", required: true)]
        public Input<string> Bucket { get; set; } = null!;

        /// <summary>
        /// Alibaba Cloud access key. It supports environment variables `ALICLOUD_ACCESS_KEY` and
        /// `ALICLOUD_ACCESS_KEY_ID`
        /// </summary>
        [Input("accessKey")]
        public Input<string>? AccessKey { get; set; }

        /// <summary>
        /// Alibaba Cloud secret access key. It supports environment variables `ALICLOUD_SECRET_KEY` and
        /// `ALICLOUD_ACCESS_KEY_SECRET`.
        /// </summary>
        [Input("secretKey")]
        public Input<string>? SecretKey { get; set; }

        /// <summary>
        /// STS access token. It supports environment variable `ALICLOUD_SECURITY_TOKEN`.
        /// </summary>
        [Input("securityToken")]
        public Input<string>? SecurityToken { get; set; }

        /// <summary>
        /// The RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the
        /// 'Access Control' section of the Alibaba Cloud console.
        /// </summary>
        [Input("ecsRoleName")]
        public Input<string>? EcsRoleName { get; set; }

        /// <summary>
        /// The region of the OSS bucket. It supports environment variables `ALICLOUD_REGION` and
        /// `ALICLOUD_DEFAULT_REGION`.
        /// </summary>
        [Input("region")]
        public Input<string>? Region { get; set; }

        /// <summary>
        /// A custom endpoint for the OSS API. It supports environment variables `ALICLOUD_OSS_ENDPOINT`
        /// and `OSS_ENDPOINT`.
        /// </summary>
        [Input("endpoint")]
        public Input<string>? Endpoint { get; set; }
        
        /// <summary>
        /// The path directory of the state file will be stored. Default to `env:`.
        /// </summary>
        [Input("prefix")]
        public Input<string>? Prefix { get; set; }

        /// <summary>
        /// The name of the state file. Defaults to `terraform.tfstate`.
        /// </summary>
        [Input("key")]
        public Input<string>? Key { get; set; }
        
        /// <summary>
        /// This is the Alibaba Cloud profile name as set in the shared credentials file. It can also be sourced from
        /// the `ALICLOUD_PROFILE` environment variable.
        /// </summary>
        [Input("profile")]
        public Input<string>? Profile { get; set; }

        /// <summary>
        /// This is the path to the shared credentials file. It can also be sourced from the
        /// `ALICLOUD_SHARED_CREDENTIALS_FILE` environment variable. If this is not set and a profile is
        /// specified, `~/.aliyun/config.json` will be used by default.
        /// </summary>
        [Input("sharedCredentialsFile")]
        public Input<string>? SharedCredentialsFile { get; set; }

        /// <summary>
        /// The ARN of the role to assume. If ARN is set to an empty string, it does not perform role switching.
        /// It supports environment variable `ALICLOUD_ASSUME_ROLE_ARN`.
        /// </summary>
        [Input("roleArn")]
        public Input<string>? RoleArn { get; set; }

        /// <summary>
        /// A more restrictive policy to apply to the temporary credentials. This gives you a way to further restrict
        /// the permissions for the resulting temporary security credentials. You cannot use this policy to grant
        /// permissions which exceed those of the role that is being assumed.
        /// </summary>
        [Input("policy")]
        public Input<string>? Policy { get; set; }

        /// <summary>
        /// The session name to use when assuming the role. It supports environment variable
        /// `ALICLOUD_ASSUME_ROLE_SESSION_NAME`
        /// </summary>
        [Input("sessionName")]
        public Input<string>? SessionName { get; set; }

        /// <summary>
        /// The time after which the established session for assuming role expires. Valid value range:
        /// [900-3600] seconds. Default to `3600`. It supports environment variable
        /// `ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION`
        /// </summary>
        [Input("sessionExpiration")]
        public Input<string>? SessionExpiration { get; set; }
    }
}