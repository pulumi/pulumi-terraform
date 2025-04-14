package main

import (
	"github.com/pulumi/pulumi-terraform/sdk/v6/go/terraform/state"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		state := state.GetLocalReferenceOutput(ctx, state.GetLocalReferenceOutputArgs{
			Path: pulumi.String("./terraform.0-12-24.tfstate"),
		}).Outputs()
		ctx.Export("state", state)
		ctx.Export("bucketArn", state.MapIndex(pulumi.String("bucket_arn")))
		ctx.Export("firstSubnetId", state.ApplyT(func(m map[string]any) string {
			return m["public_subnet_ids"].([]any)[0].(string)
		}))
		return nil
	})
}
