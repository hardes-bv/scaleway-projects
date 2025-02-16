package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"codeberg.org/hardes/iac-projects/hardes"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		org := cfg.Require("organizationId")

		err := hardes.NewDnsProject(ctx, org)
		if err != nil {
			return err
		}

		err = hardes.NewServicesProject(ctx, org)
		if err != nil {
			return err
		}
		return nil
	})
}
