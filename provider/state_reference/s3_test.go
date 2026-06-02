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
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	const (
		bucket   = "pulumi-terraform-test"
		key      = "env/terraform.tfstate"
		username = "minioadmin"
		password = "minioadmin"
	)

	container, err := tcminio.Run(ctx, "minio/minio:RELEASE.2024-12-18T13-15-44Z",
		tcminio.WithUsername(username), tcminio.WithPassword(password))
	require.NoError(t, err)

	hostPort, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	endpoint := "http://" + hostPort

	seedS3State(ctx, t, endpoint, bucket, key, username, password)

	InitTfBackend()
	resp, err := (&GetS3Reference{}).Invoke(ctx, infer.FunctionRequest[GetS3ReferenceArgs]{
		Input: GetS3ReferenceArgs{
			Workspace:                 ptr("default"),
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
