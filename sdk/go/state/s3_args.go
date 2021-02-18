package state

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"

	"github.com/pulumi/pulumi-terraform/sdk/v4/go/state/internal"
)

// S3Args specifies the configuration options for a Terraform Remote State
// stored in the S3 backend
type S3Args struct {
	// Bucket is the name of the S3 Bucket
	Bucket pulumi.StringInput

	// Key is the path to the statefile inside the bucket When using a non-default
	// workspace, the statepath will be `/workspace_key_prefix/workspace_name/key`.
	Key pulumi.StringInput

	// Regions of the S3 bucket. Sourced from `AWS_DEFAULT_REGION` environment var if unset.
	Region pulumi.StringPtrInput

	// Endpoint is a custom endpoint for the S3 API.
	// Sourced from `AWS_S3_ENDPOINT` environment var if unset.
	Endpoint pulumi.StringPtrInput

	// AccessKey is AWS Access Key ID. Sourced from the standard credentials pipeline if unset.
	AccessKey pulumi.StringPtrInput

	// SecretKey is AWS Secret Access Key. Sourced from the standard credentials pipeline if unset.
	SecretKey pulumi.StringPtrInput

	// Profile is AWS Profile name as set in the shared credentials file.
	Profile pulumi.StringPtrInput

	// SharedCredentialsFile is the path to the shared credentials file. If this is not set and a profile is
	// specified, `~/.aws/credentials` will be used by default.
	SharedCredentialsFile pulumi.StringPtrInput

	// Token is a MFA token. Sourced from `AWS_SESSION_TOKEN` environment var if unset.
	Token pulumi.StringPtrInput

	// RoleArn is an IAM Role ARN to be assumed in order to read the state from S3.
	RoleArn pulumi.StringPtrInput

	// ExternalID is the external ID to use when assuming the IAM Role
	ExternalID pulumi.StringPtrInput

	// SessionName is the session name to use when assuming the IAM Role
	SessonName pulumi.StringPtrInput

	// WorkspaceKeyPrefix is the prefix applied to the state path inside the bucket. This is only
	// relevant when using a non-default workspace and defaults to `env:`.
	WorkspaceKeyPrefix pulumi.StringPtrInput

	// IAMEndpoint is a custom endpoint for the IAM API. Sources from `AWS_IAM_ENDPOINT` if unset.
	IAMEndpoint pulumi.StringPtrInput

	// STSEndpoint is a custom endpoint for the STS API. Sources from `AWS_STS_ENDPOINT` if unset.
	STSEndpoint pulumi.StringPtrInput

	// Workspace is the Terraform workspace from which to read state
	Workspace pulumi.StringPtrInput
}

func (a *S3Args) toInternalArgs() pulumi.Input {
	return internal.S3StateReferenceArgs{
		BackendType:           pulumi.String("s3"),
		Bucket:                a.Bucket,
		Key:                   a.Key,
		Region:                a.Region,
		Endpoint:              a.Endpoint,
		AccessKey:             a.AccessKey,
		SecretKey:             a.SecretKey,
		Profile:               a.Profile,
		SharedCredentialsFile: a.SharedCredentialsFile,
		Token:                 a.Token,
		RoleArn:               a.RoleArn,
		ExternalID:            a.ExternalID,
		SessionName:           a.SessonName,
		WorkspaceKeyPrefix:    a.WorkspaceKeyPrefix,
		IAMEndpoint:           a.IAMEndpoint,
		STSEndpoint:           a.STSEndpoint,
		Workspace:             a.Workspace,
	}
}

func (l *S3Args) validateArgs() error {
	if l.Bucket == pulumi.String("") || l.Key == pulumi.String("") {
		return errors.New("`Bucket` and `Key` are required parameters")
	}
	return nil
}
