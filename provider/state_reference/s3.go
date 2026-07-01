// Copyright 2016-2025, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"

	"github.com/hashicorp/terraform/shim"
	"github.com/zclconf/go-cty/cty"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type GetS3Reference struct{}

var (
	_ = (infer.Annotated)((*GetS3Reference)(nil))
	_ = (infer.ExplicitDependencies[GetS3ReferenceArgs, StateReferenceOutputs])((*GetS3Reference)(nil))
)

func (r *GetS3Reference) Annotate(a infer.Annotator) {
	a.Describe(&r, "Access state from an AWS S3 bucket.")
}

// Taken from https://developer.hashicorp.com/terraform/language/backend/s3#configuration-variables
//
// Only the arguments that affect reading state are exposed. Write- and lock-only
// arguments (such as acl and dynamodb_table) have no effect on a read and are omitted.
type GetS3ReferenceArgs struct {
	Bucket    string  `pulumi:"bucket"`
	Key       string  `pulumi:"key"`
	Region    *string `pulumi:"region,optional"`
	Workspace *string `pulumi:"workspace,optional"`

	Endpoint    *string `pulumi:"endpoint,optional"`
	StsEndpoint *string `pulumi:"stsEndpoint,optional"`
	IamEndpoint *string `pulumi:"iamEndpoint,optional"`

	ForcePathStyle *bool `pulumi:"forcePathStyle,optional"`

	AccessKey             *string `pulumi:"accessKey,optional" provider:"secret"`
	SecretKey             *string `pulumi:"secretKey,optional" provider:"secret"`
	Token                 *string `pulumi:"token,optional" provider:"secret"`
	Profile               *string `pulumi:"profile,optional"`
	SharedCredentialsFile *string `pulumi:"sharedCredentialsFile,optional"`

	Encrypt        *bool   `pulumi:"encrypt,optional"`
	KmsKeyID       *string `pulumi:"kmsKeyId,optional"`
	SseCustomerKey *string `pulumi:"sseCustomerKey,optional" provider:"secret"`

	WorkspaceKeyPrefix *string `pulumi:"workspaceKeyPrefix,optional"`
	MaxRetries         *int    `pulumi:"maxRetries,optional"`

	SkipCredentialsValidation *bool `pulumi:"skipCredentialsValidation,optional"`
	SkipRegionValidation      *bool `pulumi:"skipRegionValidation,optional"`
	SkipMetadataAPICheck      *bool `pulumi:"skipMetadataApiCheck,optional"`
}

func (r *GetS3ReferenceArgs) Annotate(a infer.Annotator) {
	a.Describe(&r.Bucket, "The name of the S3 bucket.")
	a.Describe(&r.Key, "The path to the state file inside the bucket. When using a non-default "+
		"workspace, the state path is /workspace_key_prefix/workspace_name/key.")
	a.Describe(&r.Region, "AWS region of the S3 bucket. Falls back to the AWS_REGION or "+
		"AWS_DEFAULT_REGION environment variables when unset.")
	a.Describe(&r.Workspace, "The Terraform workspace to read state from.")

	a.Describe(&r.Endpoint, "A custom endpoint for the S3 API.")
	a.Describe(&r.StsEndpoint, "A custom endpoint for the STS API.")
	a.Describe(&r.IamEndpoint, "A custom endpoint for the IAM API.")

	a.Describe(&r.ForcePathStyle, "Force s3 to use path-style addressing instead of virtual hosted-bucket "+
		"addressing. Required by most S3-compatible stores.")

	a.Describe(&r.AccessKey, "AWS access key.")
	a.Describe(&r.SecretKey, "AWS secret key.")
	a.Describe(&r.Token, "AWS session token.")
	a.Describe(&r.Profile, "AWS profile name as set in the shared credentials file.")
	a.Describe(&r.SharedCredentialsFile, "Path to a shared credentials file.")

	a.Describe(&r.Encrypt, "Whether to enable server side encryption of the state file.")
	a.Describe(&r.KmsKeyID, "The ARN of a KMS Key to use for encrypting the state.")
	a.Describe(&r.SseCustomerKey, "The base64-encoded encryption key to use for server-side "+
		"encryption with customer-provided keys (SSE-C).")

	a.Describe(&r.WorkspaceKeyPrefix, "The prefix applied to the non-default state path inside the bucket.")
	a.Describe(&r.MaxRetries, "The maximum number of times an AWS API request is retried on retryable failure.")

	a.Describe(&r.SkipCredentialsValidation, "Skip the credentials validation via the STS API.")
	a.Describe(&r.SkipRegionValidation, "Skip static validation of region name.")
	a.Describe(&r.SkipMetadataAPICheck, "Skip the AWS Metadata API check.")

	a.SetDefault(&r.Workspace, defaultWorkspace)
}

// WireDependencies lets us tell users that our outputs shouldn't be secret, even when
// the credentials (when provided) are always secret.
//
// TODO[https://github.com/pulumi/pulumi-go-provider/issues/323]: This doesn't currently
// work; [infer.ExplicitDependencies] is not currently implemented for [infer] based
// functions.
func (r *GetS3Reference) WireDependencies(
	f infer.FieldSelector, _ *GetS3ReferenceArgs, state *StateReferenceOutputs,
) {
	f.OutputField(&state).NeverSecret() // The output should never be secret by default
}

func (r *GetS3Reference) Invoke(
	ctx context.Context, req infer.FunctionRequest[GetS3ReferenceArgs],
) (infer.FunctionResponse[StateReferenceOutputs], error) {
	args := req.Input

	results, err := shim.StateReferenceRead(ctx, "s3", *args.Workspace, map[string]cty.Value{
		"bucket":                      cty.StringVal(args.Bucket),
		"key":                         cty.StringVal(args.Key),
		"region":                      ctyStringOrNil(args.Region),
		"endpoint":                    ctyStringOrNil(args.Endpoint),
		"sts_endpoint":                ctyStringOrNil(args.StsEndpoint),
		"iam_endpoint":                ctyStringOrNil(args.IamEndpoint),
		"force_path_style":            ctyBoolOrNil(args.ForcePathStyle),
		"access_key":                  ctyStringOrNil(args.AccessKey),
		"secret_key":                  ctyStringOrNil(args.SecretKey),
		"token":                       ctyStringOrNil(args.Token),
		"profile":                     ctyStringOrNil(args.Profile),
		"shared_credentials_file":     ctyStringOrNil(args.SharedCredentialsFile),
		"encrypt":                     ctyBoolOrNil(args.Encrypt),
		"kms_key_id":                  ctyStringOrNil(args.KmsKeyID),
		"sse_customer_key":            ctyStringOrNil(args.SseCustomerKey),
		"workspace_key_prefix":        ctyStringOrNil(args.WorkspaceKeyPrefix),
		"max_retries":                 ctyIntOrNil(args.MaxRetries),
		"skip_credentials_validation": ctyBoolOrNil(args.SkipCredentialsValidation),
		"skip_region_validation":      ctyBoolOrNil(args.SkipRegionValidation),
		"skip_metadata_api_check":     ctyBoolOrNil(args.SkipMetadataAPICheck),
	})

	return infer.FunctionResponse[StateReferenceOutputs]{Output: StateReferenceOutputs{results}}, err
}
