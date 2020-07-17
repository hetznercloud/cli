package cli

import (
	"fmt"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkRemoveRouteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-route NETWORK FLAGS",
		Short:                 "Remove a route from a network",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runNetworkRemoveRoute),
	}

	cmd.Flags().IPNet("destination", net.IPNet{}, "Destination network or host (required)")
	cmd.MarkFlagRequired("destination")

	cmd.Flags().IP("gateway", net.IP{}, "Gateway IP address (required)")
	cmd.MarkFlagRequired("gateway")

	return cmd
}

func runNetworkRemoveRoute(cli *CLI, cmd *cobra.Command, args []string) error {
	gateway, _ := cmd.Flags().GetIP("gateway")
	destination, _ := cmd.Flags().GetIPNet("destination")
	idOrName := args[0]
	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
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
	action, _, err := cli.Client().Network.DeleteRoute(cli.Context, network, opts)
	if err != nil {
		return err
	}
	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Route removed from network %d\n", network.ID)

	return nil
}
