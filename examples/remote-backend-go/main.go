package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-terraform/sdk/go/state/provider"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		token := conf.RequireSecret("tfeToken")
		organization := conf.Require("organization")
		workspace := conf.Require("workspaceName")

		state, err := provider.RemoteStateReference(ctx, &provider.RemoteStateReferenceArgs{
			// TODO: This should be optional as it's hardcoded somewhere by default
			BackendType: "remote",
			BackendConfig: provider.BackendConfig{
				Organization: organization,
				Token:        token.ElementType().String(),
			},
			Workspaces: provider.WorkspaceStateArgs{
				Name: &workspace,
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("test", state.Outputs())

		return nil
	})
}
