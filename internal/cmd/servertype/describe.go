package servertype

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "serverType",
	ShortDescription:     "Describe a server type",
	JSONKeyGetByID:       "server_type",
	JSONKeyGetByName:     "server_types",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.ServerType().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.ServerType().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, _ hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		serverType := resource.(*hcloud.ServerType)

		cmd.Printf("ID:\t\t\t%d\n", serverType.ID)
		cmd.Printf("Name:\t\t\t%s\n", serverType.Name)
		cmd.Printf("Description:\t\t%s\n", serverType.Description)
		cmd.Printf("Cores:\t\t\t%d\n", serverType.Cores)
		cmd.Printf("CPU Type:\t\t%s\n", serverType.CPUType)
		cmd.Printf("Architecture:\t\t%s\n", serverType.Architecture)
		cmd.Printf("Memory:\t\t\t%.1f GB\n", serverType.Memory)
		cmd.Printf("Disk:\t\t\t%d GB\n", serverType.Disk)
		cmd.Printf("Storage Type:\t\t%s\n", serverType.StorageType)
		cmd.Printf("Included Traffic:\t%d TB\n", serverType.IncludedTraffic/util.Tebibyte)
		cmd.Printf(util.DescribeDeprecation(serverType))

		cmd.Printf("Pricings per Location:\n")
		for _, price := range serverType.Pricings {
			cmd.Printf("  - Location:\t%s:\n", price.Location.Name)
			cmd.Printf("    Hourly:\t€ %s\n", price.Hourly.Gross)
			cmd.Printf("    Monthly:\t€ %s\n", price.Monthly.Gross)
		}
		return nil
	},
}
