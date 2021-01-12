package network

import (
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newAddRouteCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-route NETWORK FLAGS",
		Short:                 "Add a route to a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.NetworkNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runAddRoute),
	}

	cmd.Flags().IPNet("destination", net.IPNet{}, "Destination network or host (required)")
	cmd.MarkFlagRequired("destination")

	cmd.Flags().IP("gateway", net.IP{}, "Gateway IP address (required)")
	cmd.MarkFlagRequired("gateway")

	return cmd
}

func runAddRoute(cli *state.State, cmd *cobra.Command, args []string) error {
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

	opts := hcloud.NetworkAddRouteOpts{
		Route: hcloud.NetworkRoute{
			Gateway:     gateway,
			Destination: &destination,
		},
	}
	action, _, err := cli.Client().Network.AddRoute(cli.Context, network, opts)
	if err != nil {
		return err
	}
	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Route added to network %d\n", network.ID)

	return nil
}
