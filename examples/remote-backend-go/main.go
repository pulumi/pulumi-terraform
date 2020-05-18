package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"

	"github.com/pulumi/pulumi-terraform/sdk/v2/go/state"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		token := conf.RequireSecret("tfeToken")
		organization := conf.Require("organization")
		workspace := conf.Require("workspaceName")

		state, err := state.NewRemoteStateReference(ctx, "remote-backend-state", &state.RemoteBackendStateArgs{
			Organization: pulumi.String(organization),
			Token:        pulumi.StringPtrInput(token.(pulumi.StringOutput)),
			Workspaces: state.WorkspaceStateArgs{
				Name: pulumi.String(workspace),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("test", state.Outputs)

		return nil
	})
}
