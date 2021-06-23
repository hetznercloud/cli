package datacenter

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "datacenter",
	ShortDescription:     "Describe an datacenter",
	JSONKeyGetByID:       "datacenter",
	JSONKeyGetByName:     "datacenters",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Datacenter().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Datacenter().Get(ctx, idOrName)
	},
	PrintText: func(ctx context.Context, client hcapi2.Client, resource interface{}) error {
		datacenter := resource.(*hcloud.Datacenter)

		fmt.Printf("ID:\t\t%d\n", datacenter.ID)
		fmt.Printf("ID:\t\t%d\n", datacenter.ID)
		fmt.Printf("Name:\t\t%s\n", datacenter.Name)
		fmt.Printf("Description:\t%s\n", datacenter.Description)
		fmt.Printf("Location:\n")
		fmt.Printf("  Name:\t\t%s\n", datacenter.Location.Name)
		fmt.Printf("  Description:\t%s\n", datacenter.Location.Description)
		fmt.Printf("  Country:\t%s\n", datacenter.Location.Country)
		fmt.Printf("  City:\t\t%s\n", datacenter.Location.City)
		fmt.Printf("  Latitude:\t%f\n", datacenter.Location.Latitude)
		fmt.Printf("  Longitude:\t%f\n", datacenter.Location.Longitude)
		fmt.Printf("Server Types:\n")

		printServerTypes := func(list []*hcloud.ServerType, client hcapi2.ServerTypeClient) {
			for _, t := range list {
				fmt.Printf("  - ID:\t\t %d\n", t.ID)
				fmt.Printf("    Name:\t %s\n", client.ServerTypeName(t.ID))
				fmt.Printf("    Description: %s\n", client.ServerTypeDescription(t.ID))
			}
		}

		fmt.Printf("  Available:\n")
		serverTypeClient := client.ServerType()
		if len(datacenter.ServerTypes.Available) > 0 {
			printServerTypes(datacenter.ServerTypes.Available, serverTypeClient)
		} else {
			fmt.Printf("    No available server types\n")
		}
		fmt.Printf("  Supported:\n")
		if len(datacenter.ServerTypes.Supported) > 0 {
			printServerTypes(datacenter.ServerTypes.Supported, serverTypeClient)
		} else {
			fmt.Printf("    No supported server types\n")
		}

		return nil
	},
}
