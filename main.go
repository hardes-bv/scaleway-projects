package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/account"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/iam"
)

func NewDnsProject(ctx *pulumi.Context, org string) error {
	orgInput := pulumi.String(org)

	// Create the Scaleway Project for `hardes.be` domain management
	dnsProjectName := "hardes"
	dnsProject, err := account.NewProject(ctx, dnsProjectName, &account.ProjectArgs{
		Name:           pulumi.String(dnsProjectName),
		Description:    pulumi.String("Hardes"),
		OrganizationId: orgInput,
	})
	if err != nil {
		return fmt.Errorf("error creating project %v: %v", dnsProjectName, err)
	}

	// Create an IAM Group
	dnsAdminGroup, err := iam.NewGroup(ctx, "hardes-dns", &iam.GroupArgs{
		OrganizationId: orgInput,
	})
	if err != nil {
		return fmt.Errorf("error creating IAM group %v: %v", dnsProjectName, err)
	}

	// Create some policies with restricted access to this project and the group as principal
	_, err = iam.NewPolicy(ctx, "hardes-dns-admin", &iam.PolicyArgs{
		GroupId:        dnsAdminGroup.ID(),
		OrganizationId: orgInput,
		Rules: iam.PolicyRuleArray{
			iam.PolicyRuleArgs{
				ProjectIds: pulumi.StringArray{
					dnsProject.ID(),
				},
				PermissionSetNames: pulumi.StringArray{
					pulumi.String("DomainsDNSFullAccess"),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating policy \"hardes-dns-admin\": %v", err)
	}

	// Create an IAM Application and API keys. Add the Application to the Group
	app, err := iam.NewApplication(ctx, "pulumi-hardes-dns", &iam.ApplicationArgs{
		OrganizationId: orgInput,
	})
	if err != nil {
		return fmt.Errorf("error creating IAM application \"pulumi-hardes-dns\": %v", err)
	}
	_, err = iam.NewApiKey(ctx, "pulumi-hardes-dns", &iam.ApiKeyArgs{
		ApplicationId: app.ID(),
	})

	// Create a Pulumi ESC environment exposing the Application API credentials

	return nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		org := cfg.Require("organizationId")

		err := NewDnsProject(ctx, org)
		if err != nil {
			return err
		}

		return nil
	})
}
