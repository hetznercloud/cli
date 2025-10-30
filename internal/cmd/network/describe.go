package network

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DescribeCmd defines a command for describing a network.
var DescribeCmd = base.DescribeCmd[*hcloud.Network]{
	ResourceNameSingular: "Network",
	ShortDescription:     "Describe a Network",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Network, any, error) {
		n, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return n, hcloud.SchemaFromNetwork(n), nil
	},
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, network *hcloud.Network) error {
		_, _ = fmt.Fprintf(out, "ID:\t%d\n", network.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", network.Name)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(network.Created), humanize.Time(network.Created))
		_, _ = fmt.Fprintf(out, "IP Range:\t%s\n", network.IPRange.String())
		_, _ = fmt.Fprintf(out, "Expose Routes to vSwitch:\t%s\n", util.YesNo(network.ExposeRoutesToVSwitch))

		if len(network.Subnets) == 0 {
			_, _ = fmt.Fprintf(out, "Subnets:\tNo subnets\n")
		} else {
			_, _ = fmt.Fprintf(out, "Subnets:\t\n")
			for _, subnet := range network.Subnets {
				_, _ = fmt.Fprintf(out, "  - Type:\t%s\n", subnet.Type)
				_, _ = fmt.Fprintf(out, "    Network Zone:\t%s\n", subnet.NetworkZone)
				_, _ = fmt.Fprintf(out, "    IP Range:\t%s\n", subnet.IPRange.String())
				_, _ = fmt.Fprintf(out, "    Gateway:\t%s\n", subnet.Gateway.String())
				if subnet.Type == hcloud.NetworkSubnetTypeVSwitch {
					_, _ = fmt.Fprintf(out, "    vSwitch ID:\t%d\n", subnet.VSwitchID)
				}
			}
		}

		if len(network.Routes) == 0 {
			_, _ = fmt.Fprintf(out, "Routes:\tNo routes\n")
		} else {
			_, _ = fmt.Fprintf(out, "Routes:\t\n")
			for _, route := range network.Routes {
				_, _ = fmt.Fprintf(out, "  - Destination:\t%s\n", route.Destination.String())
				_, _ = fmt.Fprintf(out, "    Gateway:\t%s\n", route.Gateway.String())
			}
		}

		_, _ = fmt.Fprintf(out, "Protection:\t\n")
		_, _ = fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(network.Protection.Delete))

		util.DescribeLabels(out, network.Labels, "")
		return nil
	},
}
