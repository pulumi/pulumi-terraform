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
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// TestStateReferenceReadS3 reads Terraform state from an S3-compatible store.
//
// It seeds the store the same way a user would: a MinIO container stands in for
// S3, and the real tofu CLI writes state through its own s3 backend. We then read
// that state back through GetS3Reference and confirm the outputs round-trip.
func TestStateReferenceReadS3(t *testing.T) {
	ctx := t.Context()
	endpoint := startSeededMinio(ctx, t)

	InitTfBackend()
	resp, err := (&GetS3Reference{}).Invoke(ctx, infer.FunctionRequest[GetS3ReferenceArgs]{
		Input: GetS3ReferenceArgs{
			Workspace:                 ptr(defaultWorkspace),
			Bucket:                    bucket,
			Key:                       key,
			Region:                    ptr("us-east-1"),
			Endpoint:                  ptr(endpoint),
			ForcePathStyle:            ptr(true),
			AccessKey:                 ptr(username),
			SecretKey:                 ptr(password),
			SkipCredentialsValidation: ptr(true),
			SkipRegionValidation:      ptr(true),
			SkipMetadataAPICheck:      ptr(true),
		},
	})
	require.NoError(t, err)

	assert.Equal(t, map[string]any{
		"greeting": "hello",
		"number":   float64(42),
	}, resp.Output.Outputs)
}

// TestStateReferenceReadS3AssumeRole reads state while authenticating through an
// assumed role: role_arn makes the backend call STS AssumeRole (against MinIO's
// STS API) and read the state with the returned temporary credentials.
func TestStateReferenceReadS3AssumeRole(t *testing.T) {
	ctx := t.Context()
	endpoint := startSeededMinio(ctx, t)
	configureAssumeRoleUser(ctx, t, endpoint)

	InitTfBackend()
	resp, err := (&GetS3Reference{}).Invoke(ctx, infer.FunctionRequest[GetS3ReferenceArgs]{
		Input: GetS3ReferenceArgs{
			Workspace:                 ptr(defaultWorkspace),
			Bucket:                    bucket,
			Key:                       key,
			Region:                    ptr("us-east-1"),
			Endpoint:                  ptr(endpoint),
			StsEndpoint:               ptr(endpoint),
			ForcePathStyle:            ptr(true),
			AccessKey:                 ptr(assumeRoleUsername),
			SecretKey:                 ptr(assumeRolePassword),
			RoleArn:                   ptr("arn:aws:iam::123456789012:role/terraform-state-reader"),
			SessionName:               ptr("pulumi-terraform-test"),
			SkipCredentialsValidation: ptr(true),
			SkipRegionValidation:      ptr(true),
			SkipMetadataAPICheck:      ptr(true),
		},
	})
	require.NoError(t, err)

	assert.Equal(t, map[string]any{
		"greeting": "hello",
		"number":   float64(42),
	}, resp.Output.Outputs)
}

const (
	bucket   = "pulumi-terraform-test"
	key      = "env/terraform.tfstate"
	username = "minioadmin"
	password = "minioadmin"

	assumeRoleUsername = "assume-role-user"
	assumeRolePassword = "assume-role-password"
)

// configureAssumeRoleUser creates a MinIO user whose policy allows reading the
// state only with STS credentials. MinIO sets aws:principaltype to "User" for
// static credentials and "AssumedRole" for credentials returned by AssumeRole.
func configureAssumeRoleUser(ctx context.Context, t *testing.T, endpoint string) {
	t.Helper()

	target, err := url.Parse(endpoint)
	require.NoError(t, err)
	port, err := strconv.Atoi(target.Port())
	require.NoError(t, err)

	mc, err := testcontainers.Run(ctx, "minio/mc:RELEASE.2024-11-21T17-21-54Z",
		testcontainers.WithEntrypoint("tail", "-f", "/dev/null"),
		testcontainers.WithHostPortAccess(port))
	testcontainers.CleanupContainer(t, mc, testcontainers.StopTimeout(0))
	require.NoError(t, err)

	aliasEndpoint := fmt.Sprintf("http://%s:%d", testcontainers.HostInternal, port)
	runMC := func(args ...string) {
		t.Helper()
		exitCode, output, err := mc.Exec(ctx, args)
		require.NoError(t, err)
		out, err := io.ReadAll(output)
		require.NoError(t, err)
		require.Equal(t, 0, exitCode, "%v failed:\n%s", args, out)
	}

	runMC("mc", "alias", "set", "test", aliasEndpoint, username, password)
	runMC("mc", "admin", "user", "add", "test", assumeRoleUsername, assumeRolePassword)

	policy := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": ["s3:GetObject", "s3:ListBucket"],
    "Resource": ["arn:aws:s3:::%[1]s", "arn:aws:s3:::%[1]s/*"],
    "Condition": {"StringEquals": {"aws:principaltype": "AssumedRole"}}
  }]
}`, bucket)
	require.NoError(t, mc.CopyToContainer(ctx, []byte(policy), "/tmp/assume-role-policy.json", 0o600))
	runMC("mc", "admin", "policy", "create", "test", "assume-role-only", "/tmp/assume-role-policy.json")
	runMC("mc", "admin", "policy", "attach", "test", "assume-role-only", "--user", assumeRoleUsername)

	// Prove that the source credentials cannot read the state directly. The
	// GetS3Reference invocation can succeed only by using AssumeRole credentials.
	directClient := s3.New(s3.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String(endpoint),
		UsePathStyle: true,
		Credentials:  credentials.NewStaticCredentialsProvider(assumeRoleUsername, assumeRolePassword, ""),
	})
	_, err = directClient.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	require.ErrorContains(t, err, "StatusCode: 403")
}

// startSeededMinio runs a MinIO container holding a Terraform state file with
// known outputs and returns its endpoint URL.
func startSeededMinio(ctx context.Context, t *testing.T) string {
	t.Helper()

	container, err := tcminio.Run(ctx, "minio/minio:RELEASE.2024-12-18T13-15-44Z",
		tcminio.WithUsername(username), tcminio.WithPassword(password))
	testcontainers.CleanupContainer(t, container, testcontainers.StopTimeout(0))
	require.NoError(t, err)

	hostPort, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	endpoint := "http://" + hostPort

	seedS3State(ctx, t, endpoint, bucket, key, username, password)
	return endpoint
}

// seedS3State creates the state bucket and runs tofu apply against it so the
// store holds a real Terraform state file with known outputs.
func seedS3State(ctx context.Context, t *testing.T, endpoint, bucket, key, username, password string) {
	t.Helper()

	s3Client := s3.New(s3.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String(endpoint),
		UsePathStyle: true,
		Credentials:  credentials.NewStaticCredentialsProvider(username, password, ""),
	})
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)})
	require.NoError(t, err)

	dir := t.TempDir()
	config := fmt.Sprintf(`terraform {
  backend "s3" {
    bucket                      = %q
    key                         = %q
    region                      = "us-east-1"
    access_key                  = %q
    secret_key                  = %q
    use_path_style              = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    endpoints = { s3 = %q }
  }
}

output "greeting" {
  value = "hello"
}

output "number" {
  value = 42
}
`, bucket, key, username, password, endpoint)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(config), 0o600))

	tofu := func(args ...string) {
		cmd := exec.CommandContext(ctx, "tofu", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"AWS_ACCESS_KEY_ID="+username,
			"AWS_SECRET_ACCESS_KEY="+password,
			"AWS_REGION=us-east-1",
			"AWS_EC2_METADATA_DISABLED=true",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "tofu %v failed:\n%s", args, out)
	}
	tofu("init", "-input=false")
	tofu("apply", "-auto-approve", "-input=false")
}
