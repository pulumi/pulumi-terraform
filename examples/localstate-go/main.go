package main

import (
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-terraform/sdk/go/terraform"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		fileName := conf.Require("statefile")

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		output := terraform.LocalStateReferenceOutput(ctx, terraform.LocalStateReferenceOutputArgs{
			Path: pulumi.String(filepath.Join(cwd, fileName)),
		})
		if err != nil {
			return err
		}

		ctx.Export("test", output)

		return nil
	})
}
