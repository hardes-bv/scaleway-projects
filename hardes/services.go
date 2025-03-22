package hardes

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	scw "github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway"
)

func NewServicesProject(ctx *pulumi.Context, org string) error {
	// Create the Scaleway Project for the services we will run on Kubernetes
	servicesProjectName := "hardes-services"
	_, err := scw.NewAccountProject(ctx, servicesProjectName, &scw.AccountProjectArgs{
		Name:           pulumi.String(servicesProjectName),
		Description:    pulumi.String("Hardes Application Services"),
		OrganizationId: pulumi.String(org),
	})
	if err != nil {
		return fmt.Errorf("error creating project %v: %v", servicesProjectName, err)
	}

	// Create an IAM Group

	// Create some policies with restricted access to this project and the group as principal

	// Create an IAM Application and API keys. Add the Application to the Group

	// Create a Pulumi ESC environment exposing the Application API credentials

	return nil
}
