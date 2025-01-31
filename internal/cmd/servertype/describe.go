package servertype

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.ServerType]{
	ResourceNameSingular: "serverType",
	ShortDescription:     "Describe a server type",
	JSONKeyGetByID:       "server_type",
	JSONKeyGetByName:     "server_types",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.ServerType().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.ServerType, any, error) {
		st, _, err := s.Client().ServerType().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return st, hcloud.SchemaFromServerType(st), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, serverType *hcloud.ServerType) error {
		cmd.Printf("ID:\t\t\t%d\n", serverType.ID)
		cmd.Printf("Name:\t\t\t%s\n", serverType.Name)
		cmd.Printf("Description:\t\t%s\n", serverType.Description)
		cmd.Printf("Cores:\t\t\t%d\n", serverType.Cores)
		cmd.Printf("CPU Type:\t\t%s\n", serverType.CPUType)
		cmd.Printf("Architecture:\t\t%s\n", serverType.Architecture)
		cmd.Printf("Memory:\t\t\t%.1f GB\n", serverType.Memory)
		cmd.Printf("Disk:\t\t\t%d GB\n", serverType.Disk)
		cmd.Printf("Storage Type:\t\t%s\n", serverType.StorageType)
		cmd.Print(util.DescribeDeprecation(serverType))

		pricings, err := fullPricingInfo(s, serverType)
		if err != nil {
			cmd.PrintErrf("failed to get prices for server type: %v", err)
		}

		if pricings != nil {
			cmd.Printf("Pricings per Location:\n")
			for _, price := range pricings {
				cmd.Printf("  - Location:\t\t%s\n", price.Location.Name)
				cmd.Printf("    Hourly:\t\t%s\n", util.GrossPrice(price.Hourly))
				cmd.Printf("    Monthly:\t\t%s\n", util.GrossPrice(price.Monthly))
				cmd.Printf("    Included Traffic:\t%s\n", humanize.IBytes(price.IncludedTraffic))
				cmd.Printf("    Additional Traffic:\t%s per TB\n", util.GrossPrice(price.PerTBTraffic))
				cmd.Printf("\n")
			}
		}

		return nil
	},
}

func fullPricingInfo(s state.State, serverType *hcloud.ServerType) ([]hcloud.ServerTypeLocationPricing, error) {
	pricing, _, err := s.Client().Pricing().Get(s)
	if err != nil {
		return nil, err
	}

	for _, price := range pricing.ServerTypes {
		if price.ServerType.ID == serverType.ID {
			return price.Pricings, nil
		}
	}

	return nil, nil
}
