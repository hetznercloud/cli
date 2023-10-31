package datacenter

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "datacenter",
	ShortDescription:     "Describe an datacenter",
	JSONKeyGetByID:       "datacenter",
	JSONKeyGetByName:     "datacenters",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Datacenter().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Datacenter().Get(ctx, idOrName)
	},
	PrintText: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		datacenter := resource.(*hcloud.Datacenter)

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

		printServerTypes := func(list []*hcloud.ServerType) {
			for _, t := range list {
				fmt.Printf("  - ID:\t\t %d\n", t.ID)
				fmt.Printf("    Name:\t %s\n", client.ServerType().ServerTypeName(t.ID))
				fmt.Printf("    Description: %s\n", client.ServerType().ServerTypeDescription(t.ID))
			}
		}

		fmt.Printf("  Available:\n")
		if len(datacenter.ServerTypes.Available) > 0 {
			printServerTypes(datacenter.ServerTypes.Available)
		} else {
			fmt.Printf("    No available server types\n")
		}
		fmt.Printf("  Supported:\n")
		if len(datacenter.ServerTypes.Supported) > 0 {
			printServerTypes(datacenter.ServerTypes.Supported)
		} else {
			fmt.Printf("    No supported server types\n")
		}

		return nil
	},
}
