package hardes

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	scw "github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway"
)

func NewServicesProject(ctx *pulumi.Context, org string) error {
	servicesProjectName := "hardes-services"
	_, err := scw.NewAccountProject(ctx, servicesProjectName, &scw.AccountProjectArgs{
		Name:           pulumi.String(servicesProjectName),
		Description:    pulumi.String("Hardes Application Services"),
		OrganizationId: pulumi.String(org),
	})
	if err != nil {
		return fmt.Errorf("error creating project %v: %v", servicesProjectName, err)
	}
	return nil
}
