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
		fmt.Fprintf(out, "ID:\t%d\n", network.ID)
		fmt.Fprintf(out, "Name:\t%s\n", network.Name)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(network.Created), humanize.Time(network.Created))
		fmt.Fprintf(out, "IP Range:\t%s\n", network.IPRange.String())
		fmt.Fprintf(out, "Expose Routes to vSwitch:\t%s\n", util.YesNo(network.ExposeRoutesToVSwitch))

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Subnets:\n")
		if len(network.Subnets) == 0 {
			fmt.Fprintf(out, "  No subnets\n")
		} else {
			for i, subnet := range network.Subnets {
				if i > 0 {
					fmt.Fprintln(out)
				}
				fmt.Fprintf(out, "  - Type:\t%s\n", subnet.Type)
				fmt.Fprintf(out, "    Network Zone:\t%s\n", subnet.NetworkZone)
				fmt.Fprintf(out, "    IP Range:\t%s\n", subnet.IPRange.String())
				fmt.Fprintf(out, "    Gateway:\t%s\n", subnet.Gateway.String())
				if subnet.Type == hcloud.NetworkSubnetTypeVSwitch {
					fmt.Fprintf(out, "    vSwitch ID:\t%d\n", subnet.VSwitchID)
				}
			}
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Routes:\n")
		if len(network.Routes) == 0 {
			fmt.Fprintf(out, "  No routes\n")
		} else {
			for _, route := range network.Routes {
				fmt.Fprintf(out, "  - Destination:\t%s\n", route.Destination.String())
				fmt.Fprintf(out, "    Gateway:\t%s\n", route.Gateway.String())
			}
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Protection:\n")
		fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(network.Protection.Delete))

		fmt.Fprintln(out)
		util.DescribeLabels(out, network.Labels, "")
		return nil
	},
}
