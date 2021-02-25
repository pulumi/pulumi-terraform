// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as pulumi from "@pulumi/pulumi";

/**
 * The configuration options for a Terraform Remote State stored in the AliCloud OSS backend.
 */
export interface OssRemoteStateReferenceArgs {
    /**
     * A constant describing the name of the Terraform backend, used as the discriminant
     * for the union of backend configurations.
     */
    readonly backendType: "oss";

    /**
     * Alibaba Cloud access key. It supports environment variables `ALICLOUD_ACCESS_KEY`
     * and `ALICLOUD_ACCESS_KEY_ID`
     */
    readonly accessKey?: pulumi.Input<string>;

    /**
     * Alibaba Cloud secret access key. It supports environment variables `ALICLOUD_SECRET_KEY`
     * and `ALICLOUD_ACCESS_KEY_SECRET`.
     */
    readonly secretKey?: pulumi.Input<string>;

    /**
     * STS access token. It supports environment variable `ALICLOUD_SECURITY_TOKEN`.
     */
    readonly securityToken?: pulumi.Input<string>;

    /**
     * The RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the
     * 'Access Control' section of the Alibaba Cloud console.
     */
    readonly ecsRoleName?: pulumi.Input<string>;

    /**
     * The region of the OSS bucket. It supports environment variables `ALICLOUD_REGION` and `ALICLOUD_DEFAULT_REGION`.
     */
    readonly region?: pulumi.Input<string>;

    /**
     * A custom endpoint for the OSS API. It supports environment variables `ALICLOUD_OSS_ENDPOINT` and `OSS_ENDPOINT`.
     */
    readonly endpoint?: pulumi.Input<string>;

    /**
     * The name of the OSS bucket.
     */
    readonly bucket: pulumi.Input<string>;

    /**
     * The path directory of the state file will be stored. Default to `env:`.
     */
    readonly prefix?: pulumi.Input<string>;

    /**
     * The name of the state file. Defaults to `terraform.tfstate`.
     */
    readonly key?: pulumi.Input<string>;

    /**
     * This is the Alibaba Cloud profile name as set in the shared credentials file. It can also be sourced from
     * the `ALICLOUD_PROFILE` environment variable.
     */
    readonly profile?: pulumi.Input<string>;

    /**
     * This is the path to the shared credentials file. It can also be sourced from the
     * `ALICLOUD_SHARED_CREDENTIALS_FILE` environment variable. If this is not set and a profile is
     * specified, `~/.aliyun/config.json` will be used.
     */
    readonly sharedCredentialsFile?: pulumi.Input<string>;

    /**
     * The ARN of the role to assume. If ARN is set to an empty string, it does not perform role switching.
     * It supports environment variable `ALICLOUD_ASSUME_ROLE_ARN`.
     */
    readonly roleArn?: pulumi.Input<string>;

    /**
     *  A more restrictive policy to apply to the temporary credentials. This gives you a way to further restrict
     *  the permissions for the resulting temporary security credentials. You cannot use this policy to grant
     *  permissions which exceed those of the role that is being assumed.
     */
    readonly policy?: pulumi.Input<string>;

    /**
     * The session name to use when assuming the role. It supports environment variable
     * `ALICLOUD_ASSUME_ROLE_SESSION_NAME`
     */
    readonly sessionName?: pulumi.Input<string>;

    /**
     * The time after which the established session for assuming role expires.
     * Valid value range: [900-3600] seconds. Default to `3600`. It supports environment variable
     * `ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION`.
     */
    readonly sessionExpiration?: pulumi.Input<string>;
}

