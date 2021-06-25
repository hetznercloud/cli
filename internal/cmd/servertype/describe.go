package servertype

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "serverType",
	ShortDescription:     "Describe a server type",
	JSONKeyGetByID:       "server_type",
	JSONKeyGetByName:     "server_types",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Image().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.ServerType().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, _ hcapi2.Client, _ *cobra.Command, resource interface{}) error {
		serverType := resource.(*hcloud.ServerType)

		fmt.Printf("ID:\t\t%d\n", serverType.ID)
		fmt.Printf("Name:\t\t%s\n", serverType.Name)
		fmt.Printf("Description:\t%s\n", serverType.Description)
		fmt.Printf("Cores:\t\t%d\n", serverType.Cores)
		fmt.Printf("CPU Type:\t%s\n", serverType.CPUType)
		fmt.Printf("Memory:\t\t%.1f GB\n", serverType.Memory)
		fmt.Printf("Disk:\t\t%d GB\n", serverType.Disk)
		fmt.Printf("Storage Type:\t%s\n", serverType.StorageType)

		fmt.Printf("Pricings per Location:\n")
		for _, price := range serverType.Pricings {
			fmt.Printf("  - Location:\t%s:\n", price.Location.Name)
			fmt.Printf("    Hourly:\t€ %s\n", price.Hourly.Gross)
			fmt.Printf("    Monthly:\t€ %s\n", price.Monthly.Gross)
		}
		return nil
	},
}
