package main

import (
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi-terraform/sdk/v3/go/state"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		fileName := conf.Require("statefile")

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		state, err := state.NewRemoteStateReference(ctx, "localstate", &state.LocalStateArgs{
			Path: pulumi.String(filepath.Join(cwd, fileName)),
		})
		if err != nil {
			return err
		}

		ctx.Export("test", state.Outputs)

		return nil
	})
}
