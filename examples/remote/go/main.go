package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-terraform/sdk/v6/go/terraform"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		token := conf.RequireSecret("tfeToken")
		organization := conf.Require("organization")
		workspacesPrefiix := conf.Require("workspacesPrefix")

		output := terraform.RemoteStateReferenceOutput(ctx, terraform.RemoteStateReferenceOutputArgs{
			Organization: pulumi.String(organization),
			Token:        pulumi.StringInput(token),
			Workspaces: terraform.WorkspacesArgs{
				Prefix: pulumi.StringPtr(workspacesPrefiix),
			},
		})
		plainOutput := pulumi.Unsecret(output)

		ctx.Export("test", plainOutput)

		return nil
	})
}
