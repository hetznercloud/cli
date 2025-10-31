package servertype

import (
	"fmt"
	"io"
	"strings"

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
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, serverType *hcloud.ServerType) error {
		description, err := DescribeServerType(s, serverType, false)
		if err != nil {
			return err
		}
		fmt.Fprint(out, description)
		return nil
	},
}

func DescribeServerType(s state.State, serverType *hcloud.ServerType, short bool) (string, error) {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ID:\t%d\n", serverType.ID)
	fmt.Fprintf(&sb, "Name:\t%s\n", serverType.Name)
	fmt.Fprintf(&sb, "Description:\t%s\n", serverType.Description)
	fmt.Fprintf(&sb, "Category:\t%s\n", serverType.Category)
	fmt.Fprintf(&sb, "Cores:\t%d\n", serverType.Cores)
	fmt.Fprintf(&sb, "CPU Type:\t%s\n", serverType.CPUType)
	fmt.Fprintf(&sb, "Architecture:\t%s\n", serverType.Architecture)
	fmt.Fprintf(&sb, "Memory:\t%.1f GB\n", serverType.Memory)
	fmt.Fprintf(&sb, "Disk:\t%d GB\n", serverType.Disk)
	fmt.Fprintf(&sb, "Storage Type:\t%s\n", serverType.StorageType)

	if short {
		return sb.String(), nil
	}

	pricings, err := fullPricingInfo(s, serverType)
	if err != nil {
		return "", fmt.Errorf("failed to get prices for Server Type: %w", err)
	}

	locations := joinLocationInfo(serverType, pricings)
	fmt.Fprintln(&sb)
	fmt.Fprintf(&sb, "Locations:\n")
	for _, info := range locations {

		fmt.Fprintf(&sb, "  - Location:\t%s\n", info.Location.Name)

		if deprecationText := util.DescribeDeprecation(info); deprecationText != "" {
			fmt.Fprint(&sb, util.PrefixLines(deprecationText, "    "))
		}

		fmt.Fprintf(&sb, "    Pricing:\t\n")
		fmt.Fprintf(&sb, "      Hourly:\t%s\n", util.GrossPrice(info.Pricing.Hourly))
		fmt.Fprintf(&sb, "      Monthly:\t%s\n", util.GrossPrice(info.Pricing.Monthly))
		fmt.Fprintf(&sb, "      Included Traffic:\t%s\n", humanize.IBytes(info.Pricing.IncludedTraffic))
		fmt.Fprintf(&sb, "      Additional Traffic:\t%s per TB\n", util.GrossPrice(info.Pricing.PerTBTraffic))
		fmt.Fprintf(&sb, "\n")
	}

	return sb.String(), nil
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
