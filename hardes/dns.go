package hardes

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	scw "github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway"
)

func NewDnsProject(ctx *pulumi.Context, org string) error {
	dnsProjectName := "hardes-dns"
	_, err := scw.NewAccountProject(ctx, dnsProjectName, &scw.AccountProjectArgs{
		Name:           pulumi.String(dnsProjectName),
		Description:    pulumi.String("Hardes Domain Name Management"),
		OrganizationId: pulumi.String(org),
	})
	if err != nil {
		return fmt.Errorf("error creating project %v: %v", dnsProjectName, err)
	}
	return nil
}
