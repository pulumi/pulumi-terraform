package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v5/go/state/internal"
)

// OssArgs specifies the configuration options for a Terraform Remote State
// stored in the Oss backend
type OssArgs struct {
	// Bucket is the name of the OSS bucket.
	Bucket pulumi.StringInput

	// AccessKey is the Alibaba Cloud access key. It supports environment variables `ALICLOUD_ACCESS_KEY` and
	// `ALICLOUD_ACCESS_KEY_ID`
	AccessKey pulumi.StringPtrInput

	// SecurityKey is the Alibaba Cloud secret access key. It supports environment variables `ALICLOUD_SECRET_KEY`
	// and `ALICLOUD_ACCESS_KEY_SECRET`.
	SecretKey pulumi.StringPtrInput

	// SecurityToken is the STS access token. It supports environment variable `ALICLOUD_SECURITY_TOKEN`.
	SecurityToken pulumi.StringPtrInput

	// EcsRoleName is the RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the
	// 'Access Control' section of the Alibaba Cloud console.
	EcsRoleName pulumi.StringPtrInput

	// Region is the region of the OSS bucket. It supports environment variables `ALICLOUD_REGION` and
	//`ALICLOUD_DEFAULT_REGION`.
	Region pulumi.StringPtrInput

	// Endpoint is a custom endpoint for the OSS API. It supports environment variables `ALICLOUD_OSS_ENDPOINT`
	// and `OSS_ENDPOINT`.
	Endpoint pulumi.StringPtrInput

	// Prefix is the path directory of the state file will be stored. Default to `env:`.
	Prefix pulumi.StringPtrInput

	// Key is the name of the state file. Defaults to `terraform.tfstate`.
	Key pulumi.StringPtrInput

	// Profile is the Alibaba Cloud profile name as set in the shared credentials file. It can also be sourced from
	// the `ALICLOUD_PROFILE` environment variable.
	Profile pulumi.StringPtrInput

	// SharedCredentialsFile is the path to the shared credentials file. It can also be sourced from the
	// `ALICLOUD_SHARED_CREDENTIALS_FILE` environment variable. If this is not set and a profile is
	// specified, `~/.aliyun/config.json` will be used.
	SharedCredentialsFile pulumi.StringPtrInput

	// RoleArn is the ARN of the role to assume. If ARN is set to an empty string, it does not perform role switching.
	RoleArn pulumi.StringPtrInput

	// Policy is a more restrictive policy to apply to the temporary credentials. This gives you a way to further
	// restrict the permissions for the resulting temporary security credentials. You cannot use this policy to grant
	// permissions which exceed those of the role that is being assumed.
	Policy pulumi.StringPtrInput

	// SessionName is the session name to use when assuming the role. It supports environment variable
	// `ALICLOUD_ASSUME_ROLE_SESSION_NAME`
	SessonName pulumi.StringPtrInput

	// SessionExpiration is the time after which the established session for assuming role expires. Valid value range:
	// [900-3600] seconds. Default to `3600`. It supports environment variable `ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION`
	SessionExpiration pulumi.StringPtrInput
}

func (a *OssArgs) toInternalArgs() pulumi.Input {
	return internal.OssStateReferenceArgs{
		BackendType:           pulumi.String("oss"),
		Bucket:                a.Bucket,
		AccessKey:             a.AccessKey,
		SecretKey:             a.SecretKey,
		SecurityToken:         a.SecurityToken,
		EcsRoleName:           a.EcsRoleName,
		Key:                   a.Key,
		Prefix:                a.Prefix,
		Region:                a.Region,
		Endpoint:              a.Endpoint,
		Profile:               a.Profile,
		SharedCredentialsFile: a.SharedCredentialsFile,
		RoleArn:               a.RoleArn,
		SessionName:           a.SessonName,
		SessionExpiration:     a.SessionExpiration,
		Policy:                a.Policy,
	}
}

func (l *OssArgs) validateArgs() error {
	if l.Bucket == pulumi.String("") {
		return errors.New("`Bucket` is a required parameter")
	}
	return nil
}
