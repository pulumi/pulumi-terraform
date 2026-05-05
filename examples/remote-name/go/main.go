package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-terraform/sdk/v6/go/terraform/state"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		token := conf.RequireSecret("remote_tf_token")
		organization := conf.Require("remote_tf_org")
		workspaceName := conf.Require("workspace_name")

		state := state.GetRemoteReferenceOutput(ctx, state.GetRemoteReferenceOutputArgs{
			Organization: pulumi.String(organization),
			Token:        pulumi.StringInput(token),
			Workspaces: state.WorkspacesArgs{
				Name: pulumi.StringPtr(workspaceName),
			},
		})
		ctx.Export("state", state.Outputs())
		return nil
	})
}
