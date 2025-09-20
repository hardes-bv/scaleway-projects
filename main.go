package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-pulumiservice/sdk/go/pulumiservice"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/account"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/iam"
)

func NewDnsProject(ctx *pulumi.Context, org string) error {
	orgInput := pulumi.String(org)

	// Create the Scaleway Project for `hardes.be` domain management
	projectName := "hardes"
	project, err := account.NewProject(ctx, projectName, &account.ProjectArgs{
		Name:           pulumi.String(projectName),
		Description:    pulumi.String("Hardes"),
		OrganizationId: orgInput,
	})
	if err != nil {
		return fmt.Errorf("error creating project %v: %v", projectName, err)
	}

	// Create an IAM Group
	adminGroup, err := iam.NewGroup(ctx, "hardes-dns", &iam.GroupArgs{
		OrganizationId: orgInput,
	})
	if err != nil {
		return fmt.Errorf("error creating IAM group %v: %v", projectName, err)
	}

	// Create some policies with restricted access to this project and the group as principal
	_, err = iam.NewPolicy(ctx, "hardes-admin", &iam.PolicyArgs{
		Name:           pulumi.String("hardes-admin"),
		GroupId:        adminGroup.ID(),
		OrganizationId: orgInput,
		Rules: iam.PolicyRuleArray{
			iam.PolicyRuleArgs{
				OrganizationId: orgInput,
				PermissionSetNames: pulumi.StringArray{
					pulumi.String("IAMManager"),
				},
			},
			iam.PolicyRuleArgs{
				ProjectIds: pulumi.StringArray{
					project.ID(),
				},
				PermissionSetNames: pulumi.StringArray{
					pulumi.String("AllProductsFullAccess"),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating policy \"hardes-admin\": %v", err)
	}

	// Create an IAM Application and API keys. Add the Application to the Group
	app, err := iam.NewApplication(ctx, "pulumi-hardes", &iam.ApplicationArgs{
		OrganizationId: orgInput,
	})
	if err != nil {
		return fmt.Errorf("error creating IAM application \"pulumi-hardes\": %v", err)
	}

	_, err = iam.NewGroupMembership(ctx, "pulumi-hardes-member", &iam.GroupMembershipArgs{
		GroupId:       adminGroup.ID(),
		ApplicationId: app.ID(),
	})
	if err != nil {
		return fmt.Errorf("error creating group membership \"pulumi-hardes\": %v", err)
	}

	apiKey, err := iam.NewApiKey(ctx, "pulumi-hardes", &iam.ApiKeyArgs{
		ApplicationId: app.ID(),
	})
	if err != nil {
		return fmt.Errorf("error creating api key \"pulumi-hardes\": %v", err)
	}

	// Create a Pulumi ESC environment exposing the Application API credentials
	_, err = pulumiservice.NewEnvironment(ctx, "environmentResource", &pulumiservice.EnvironmentArgs{
		Name:         pulumi.String("hardes"),
		Organization: pulumi.String(ctx.Organization()),
		Yaml:         pulumi.NewFileAsset("environment.yaml"),
		Project:      pulumi.String("scaleway"),
	})
	if err != nil {
		return fmt.Errorf("error creating ESC environment \"scaleway\\hardes\": %v", err)
	}

	// Stack exports
	ctx.Export("organizationId", orgInput)
	ctx.Export("projectId", project.ID())
	ctx.Export("accessKey", apiKey.AccessKey)
	ctx.Export("secretKey", apiKey.SecretKey)

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
