package main

import (
	"github.com/pulumi/pulumi-terraform/sdk/v6/go/terraform"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		state := terraform.LocalStateReferenceOutput(ctx, terraform.LocalStateReferenceOutputArgs{
			Path: pulumi.String("./terraform.0-12-24.tfstate"),
		})
		ctx.Export("state", state.Outputs())
		return nil
	})
}
