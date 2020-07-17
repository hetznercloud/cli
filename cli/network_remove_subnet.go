package cli

import (
	"fmt"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkRemoveSubnetCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-subnet NETWORK FLAGS",
		Short:                 "Remove a subnet from a network",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runNetworkRemoveSubnet),
	}
	cmd.Flags().IPNet("ip-range", net.IPNet{}, "Subnet IP range (required)")
	cmd.MarkFlagRequired("ip-range")
	return cmd
}

func runNetworkRemoveSubnet(cli *CLI, cmd *cobra.Command, args []string) error {

	ipRange, _ := cmd.Flags().GetIPNet("ip-range")
	idOrName := args[0]
	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}

	opts := hcloud.NetworkDeleteSubnetOpts{
		Subnet: hcloud.NetworkSubnet{
			IPRange: &ipRange,
		},
	}
	action, _, err := cli.Client().Network.DeleteSubnet(cli.Context, network, opts)
	if err != nil {
		return err
	}
	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Subnet %s removed from network %d\n", ipRange.String(), network.ID)

	return nil
}
