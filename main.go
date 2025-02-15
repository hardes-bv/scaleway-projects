package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	scw "github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		org := cfg.Require("organizationId")

		_, err := scw.NewAccountProject(ctx, "hardes-dns", &scw.AccountProjectArgs{
			Name:           pulumi.String("hardes-dns"),
			Description:    pulumi.String("Hardes Domain Name Management"),
			OrganizationId: pulumi.String(org),
		})
		if err != nil {
			return fmt.Errorf("error creating public IP: %v", err)
		}
		return nil
	})
}
