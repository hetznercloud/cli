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
	ResourceNameSingular: "Server Type",
	ShortDescription:     "Describe a Server Type",
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
		cmd.Printf("Category:\t\t%s\n", serverType.Category)
		cmd.Printf("Cores:\t\t\t%d\n", serverType.Cores)
		cmd.Printf("CPU Type:\t\t%s\n", serverType.CPUType)
		cmd.Printf("Architecture:\t\t%s\n", serverType.Architecture)
		cmd.Printf("Memory:\t\t\t%.1f GB\n", serverType.Memory)
		cmd.Printf("Disk:\t\t\t%d GB\n", serverType.Disk)
		cmd.Printf("Storage Type:\t\t%s\n", serverType.StorageType)

		pricings, err := fullPricingInfo(s, serverType)
		if err != nil {
			cmd.PrintErrf("failed to get prices for Server Type: %v", err)
		}

		locations := joinLocationInfo(serverType, pricings)
		cmd.Printf("Locations:\n")
		for _, info := range locations {

			cmd.Printf("  - Location:\t\t\t%s\n", info.Location.Name)

			if deprecationText := util.DescribeDeprecation(info); deprecationText != "" {
				cmd.Print(util.PrefixLines(deprecationText, "    "))
			}

			cmd.Printf("    Pricing:\n")
			cmd.Printf("      Hourly:\t\t\t%s\n", util.GrossPrice(info.Pricing.Hourly))
			cmd.Printf("      Monthly:\t\t\t%s\n", util.GrossPrice(info.Pricing.Monthly))
			cmd.Printf("      Included Traffic:\t\t%s\n", humanize.IBytes(info.Pricing.IncludedTraffic))
			cmd.Printf("      Additional Traffic:\t%s per TB\n", util.GrossPrice(info.Pricing.PerTBTraffic))
			cmd.Printf("\n")
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

type locationInfo struct {
	Location *hcloud.Location
	hcloud.DeprecatableResource
	Pricing hcloud.ServerTypeLocationPricing
}

func joinLocationInfo(serverType *hcloud.ServerType, pricings []hcloud.ServerTypeLocationPricing) []locationInfo {
	locations := make([]locationInfo, 0, len(serverType.Locations))

	for _, location := range serverType.Locations {
		info := locationInfo{Location: location.Location, DeprecatableResource: location.DeprecatableResource}

		for _, pricing := range pricings {
			// Pricing endpoint only sets the location name
			if pricing.Location.Name == info.Location.Name {
				info.Pricing = pricing
				break
			}
		}

		locations = append(locations, info)
	}

	return locations
}
