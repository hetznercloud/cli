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
		_, _ = fmt.Fprint(out, description)
		return nil
	},
}

func DescribeServerType(s state.State, serverType *hcloud.ServerType, short bool) (string, error) {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "ID:\t%d\n", serverType.ID)
	_, _ = fmt.Fprintf(&sb, "Name:\t%s\n", serverType.Name)
	_, _ = fmt.Fprintf(&sb, "Description:\t%s\n", serverType.Description)
	_, _ = fmt.Fprintf(&sb, "Category:\t%s\n", serverType.Category)
	_, _ = fmt.Fprintf(&sb, "Cores:\t%d\n", serverType.Cores)
	_, _ = fmt.Fprintf(&sb, "CPU Type:\t%s\n", serverType.CPUType)
	_, _ = fmt.Fprintf(&sb, "Architecture:\t%s\n", serverType.Architecture)
	_, _ = fmt.Fprintf(&sb, "Memory:\t%.1f GB\n", serverType.Memory)
	_, _ = fmt.Fprintf(&sb, "Disk:\t%d GB\n", serverType.Disk)
	_, _ = fmt.Fprintf(&sb, "Storage Type:\t%s\n", serverType.StorageType)

	if short {
		return sb.String(), nil
	}

	pricings, err := fullPricingInfo(s, serverType)
	if err != nil {
		return "", fmt.Errorf("failed to get prices for Server Type: %v", err)
	}

	locations := joinLocationInfo(serverType, pricings)
	_, _ = fmt.Fprintf(&sb, "Locations:\n")
	for _, info := range locations {

		_, _ = fmt.Fprintf(&sb, "  - Location:\t\t\t%s\n", info.Location.Name)

		if deprecationText := util.DescribeDeprecation(info); deprecationText != "" {
			_, _ = fmt.Fprintf(&sb, util.PrefixLines(deprecationText, "    "))
		}

		_, _ = fmt.Fprintf(&sb, "    Pricing:\n")
		_, _ = fmt.Fprintf(&sb, "      Hourly:\t%s\n", util.GrossPrice(info.Pricing.Hourly))
		_, _ = fmt.Fprintf(&sb, "      Monthly:\t%s\n", util.GrossPrice(info.Pricing.Monthly))
		_, _ = fmt.Fprintf(&sb, "      Included Traffic:\t%s\n", humanize.IBytes(info.Pricing.IncludedTraffic))
		_, _ = fmt.Fprintf(&sb, "      Additional Traffic:\t%s per TB\n", util.GrossPrice(info.Pricing.PerTBTraffic))
		_, _ = fmt.Fprintf(&sb, "\n")
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
