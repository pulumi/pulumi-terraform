package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-terraform/sdk/go/terraform"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		token := conf.RequireSecret("tfeToken")
		organization := conf.Require("organization")
		workspace := conf.Require("workspaceName")

		output := terraform.RemoteStateReferenceOutput(ctx, terraform.RemoteStateReferenceOutputArgs{
			BackendConfig: terraform.BackendConfigArgs{
				Organization: pulumi.String(organization),
				Token:        pulumi.StringInput(token),
			},
			Workspaces: terraform.WorkspaceArgs{
				Name: pulumi.StringPtr(workspace),
			},
		})

		ctx.Export("test", output)

		return nil
	})
}
