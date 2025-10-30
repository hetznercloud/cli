package loadbalancertype

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

var DescribeCmd = base.DescribeCmd[*hcloud.LoadBalancerType]{
	ResourceNameSingular: "Load Balancer Type",
	ShortDescription:     "Describe a Load Balancer Type",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancerType().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.LoadBalancerType, any, error) {
		lbt, _, err := s.Client().LoadBalancerType().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return lbt, hcloud.SchemaFromLoadBalancerType(lbt), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, loadBalancerType *hcloud.LoadBalancerType) error {
		_, _ = fmt.Fprint(out, DescribeLoadBalancerType(s, loadBalancerType, false))
		return nil
	},
}

func DescribeLoadBalancerType(s state.State, loadBalancerType *hcloud.LoadBalancerType, short bool) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "ID:\t%d\n", loadBalancerType.ID)
	_, _ = fmt.Fprintf(&sb, "Name:\t%s\n", loadBalancerType.Name)
	_, _ = fmt.Fprintf(&sb, "Description:\t%s\n", loadBalancerType.Description)
	_, _ = fmt.Fprintf(&sb, "Max Services:\t%d\n", loadBalancerType.MaxServices)
	_, _ = fmt.Fprintf(&sb, "Max Connections:\t%d\n", loadBalancerType.MaxConnections)
	_, _ = fmt.Fprintf(&sb, "Max Targets:\t%d\n", loadBalancerType.MaxTargets)
	_, _ = fmt.Fprintf(&sb, "Max assigned Certificates:\t%d\n", loadBalancerType.MaxAssignedCertificates)

	if short {
		return sb.String()
	}

	pricings, err := fullPricingInfo(s, loadBalancerType)
	if err != nil {
		_, _ = fmt.Fprintf(&sb, "failed to get prices for Load Balancer Type: %v", err)
	}

	if pricings != nil {
		_, _ = fmt.Fprintf(&sb, "Pricings per Location:\t\n")
		for _, price := range pricings {
			_, _ = fmt.Fprintf(&sb, "  - Location:\t%s\n", price.Location.Name)
			_, _ = fmt.Fprintf(&sb, "    Hourly:\t%s\n", util.GrossPrice(price.Hourly))
			_, _ = fmt.Fprintf(&sb, "    Monthly:\t%s\n", util.GrossPrice(price.Monthly))
			_, _ = fmt.Fprintf(&sb, "    Included Traffic:\t%s\n", humanize.IBytes(price.IncludedTraffic))
			_, _ = fmt.Fprintf(&sb, "    Additional Traffic:\t%s per TB\n", util.GrossPrice(price.PerTBTraffic))
			_, _ = fmt.Fprintf(&sb, "\t\n")
		}
	}

	return sb.String()
}

func fullPricingInfo(s state.State, loadBalancerType *hcloud.LoadBalancerType) ([]hcloud.LoadBalancerTypeLocationPricing, error) {
	pricing, _, err := s.Client().Pricing().Get(s)
	if err != nil {
		return nil, err
	}

	for _, price := range pricing.LoadBalancerTypes {
		if price.LoadBalancerType.ID == loadBalancerType.ID {
			return price.Pricings, nil
		}
	}

	return nil, nil
}
