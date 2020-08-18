package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"

	"github.com/pulumi/pulumi-terraform/sdk/v3/go/state"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		bucket := conf.Require("bucketName")
		key := conf.Require("key")
		region := conf.Require("region")

		state, err := state.NewRemoteStateReference(ctx, "s3state", &state.S3Args{
			Bucket: pulumi.String(bucket),
			Key:    pulumi.Sprintf("%s/terraform.tfstate", key),
			Region: pulumi.String(region),
		})
		if err != nil {
			return err
		}

		ctx.Export("test", state.Outputs)

		return nil
	})
}
