package hardes

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/account"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/iam"
)

func NewDnsProject(ctx *pulumi.Context, org string) error {
	// Create the Scaleway Project for `hardes.be` domain management
	dnsProjectName := "hardes-dns"
	_, err := account.NewProject(ctx, dnsProjectName, &account.ProjectArgs{
		Name:           pulumi.String(dnsProjectName),
		Description:    pulumi.String("Hardes Domain Name Management"),
		OrganizationId: pulumi.String(org),
	})
	if err != nil {
		return fmt.Errorf("error creating project %v: %v", dnsProjectName, err)
	}

	// Create an IAM Group
	_, err = iam.NewGroup(ctx, "hardes-dns", &iam.GroupArgs{
		OrganizationId: pulumi.String(org),
	})
	if err != nil {
		return fmt.Errorf("error creating IAM group %v: %v", dnsProjectName, err)
	}

	// Create some policies with restricted access to this project and the group as principal

	// Create an IAM Application and API keys. Add the Application to the Group

	// Create a Pulumi ESC environment exposing the Application API credentials

	return nil
}
