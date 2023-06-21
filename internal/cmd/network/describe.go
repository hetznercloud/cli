package network

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

// DescribeCmd defines a command for describing a network.
var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "network",
	ShortDescription:     "Describe a network",
	JSONKeyGetByID:       "network",
	JSONKeyGetByName:     "networks",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Network().Get(ctx, idOrName)
	},
	PrintText: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		network := resource.(*hcloud.Network)

		fmt.Printf("ID:\t\t%d\n", network.ID)
		fmt.Printf("Name:\t\t%s\n", network.Name)
		fmt.Printf("Created:\t%s (%s)\n", util.Datetime(network.Created), humanize.Time(network.Created))
		fmt.Printf("IP Range:\t%s\n", network.IPRange.String())
		fmt.Printf("Expose Routes to vSwitch: %s\n", util.YesNo(network.ExposeRoutesToVSwitch))

		fmt.Printf("Subnets:\n")
		if len(network.Subnets) == 0 {
			fmt.Print("  No subnets\n")
		} else {
			for _, subnet := range network.Subnets {
				fmt.Printf("  - Type:\t\t%s\n", subnet.Type)
				fmt.Printf("    Network Zone:\t%s\n", subnet.NetworkZone)
				fmt.Printf("    IP Range:\t\t%s\n", subnet.IPRange.String())
				fmt.Printf("    Gateway:\t\t%s\n", subnet.Gateway.String())
				if subnet.Type == hcloud.NetworkSubnetTypeVSwitch {
					fmt.Printf("    vSwitch ID:\t\t%d\n", subnet.VSwitchID)
				}
			}
		}

		fmt.Printf("Routes:\n")
		if len(network.Routes) == 0 {
			fmt.Print("  No routes\n")
		} else {
			for _, route := range network.Routes {
				fmt.Printf("  - Destination:\t%s\n", route.Destination.String())
				fmt.Printf("    Gateway:\t\t%s\n", route.Gateway.String())
			}
		}

		fmt.Printf("Protection:\n")
		fmt.Printf("  Delete:\t%s\n", util.YesNo(network.Protection.Delete))

		fmt.Print("Labels:\n")
		if len(network.Labels) == 0 {
			fmt.Print("  No labels\n")
		} else {
			for key, value := range network.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
