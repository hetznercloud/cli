package network

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DescribeCmd defines a command for describing a network.
var DescribeCmd = base.DescribeCmd[*hcloud.Network]{
	ResourceNameSingular: "network",
	ShortDescription:     "Describe a network",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Network, any, error) {
		n, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return n, hcloud.SchemaFromNetwork(n), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, network *hcloud.Network) error {
		cmd.Printf("ID:\t\t%d\n", network.ID)
		cmd.Printf("Name:\t\t%s\n", network.Name)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(network.Created), humanize.Time(network.Created))
		cmd.Printf("IP Range:\t%s\n", network.IPRange.String())
		cmd.Printf("Expose Routes to vSwitch: %s\n", util.YesNo(network.ExposeRoutesToVSwitch))

		cmd.Printf("Subnets:\n")
		if len(network.Subnets) == 0 {
			cmd.Print("  No subnets\n")
		} else {
			for _, subnet := range network.Subnets {
				cmd.Printf("  - Type:\t\t%s\n", subnet.Type)
				cmd.Printf("    Network Zone:\t%s\n", subnet.NetworkZone)
				cmd.Printf("    IP Range:\t\t%s\n", subnet.IPRange.String())
				cmd.Printf("    Gateway:\t\t%s\n", subnet.Gateway.String())
				if subnet.Type == hcloud.NetworkSubnetTypeVSwitch {
					cmd.Printf("    vSwitch ID:\t\t%d\n", subnet.VSwitchID)
				}
			}
		}

		cmd.Printf("Routes:\n")
		if len(network.Routes) == 0 {
			cmd.Print("  No routes\n")
		} else {
			for _, route := range network.Routes {
				cmd.Printf("  - Destination:\t%s\n", route.Destination.String())
				cmd.Printf("    Gateway:\t\t%s\n", route.Gateway.String())
			}
		}

		cmd.Printf("Protection:\n")
		cmd.Printf("  Delete:\t%s\n", util.YesNo(network.Protection.Delete))

		cmd.Print("Labels:\n")
		if len(network.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range network.Labels {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
