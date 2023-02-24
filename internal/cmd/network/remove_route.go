package network

import (
	"context"
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var RemoveRouteCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "remove-route NETWORK FLAGS",
			Short:                 "Remove a route from a network",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().IPNet("destination", net.IPNet{}, "Destination network or host (required)")
		cmd.MarkFlagRequired("destination")

		cmd.Flags().IP("gateway", net.IP{}, "Gateway IP address (required)")
		cmd.MarkFlagRequired("gateway")

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		gateway, _ := cmd.Flags().GetIP("gateway")
		destination, _ := cmd.Flags().GetIPNet("destination")
		idOrName := args[0]
		network, _, err := client.Network().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", idOrName)
		}

		opts := hcloud.NetworkDeleteRouteOpts{
			Route: hcloud.NetworkRoute{
				Gateway:     gateway,
				Destination: &destination,
			},
		}
		action, _, err := client.Network().DeleteRoute(ctx, network, opts)
		if err != nil {
			return err
		}
		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}
		fmt.Printf("Route removed from network %d\n", network.ID)

		return nil
	},
}
