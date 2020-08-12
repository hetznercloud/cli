package cli

import (
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkAddSubnetCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-subnet NETWORK FLAGS",
		Short:                 "Add a subnet to a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.NetworkNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runNetworkAddSubnet),
	}

	cmd.Flags().String("type", "", "Type of subnet")
	cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("cloud", "server"))
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("network-zone", "", "Name of network zone")
	cmd.RegisterFlagCompletionFunc("network-zone", cmpl.SuggestCandidates("eu-central"))
	cmd.MarkFlagRequired("network-zone")

	cmd.Flags().IPNet("ip-range", net.IPNet{}, "Range to allocate IPs from")

	return cmd
}

func runNetworkAddSubnet(cli *CLI, cmd *cobra.Command, args []string) error {
	subnetType, _ := cmd.Flags().GetString("type")
	networkZone, _ := cmd.Flags().GetString("network-zone")
	ipRange, _ := cmd.Flags().GetIPNet("ip-range")
	idOrName := args[0]

	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}
	subnet := hcloud.NetworkSubnet{
		Type:        hcloud.NetworkSubnetType(subnetType),
		NetworkZone: hcloud.NetworkZone(networkZone),
	}

	if ipRange.IP != nil && ipRange.Mask != nil {
		subnet.IPRange = &ipRange
	}

	opts := hcloud.NetworkAddSubnetOpts{
		Subnet: subnet,
	}
	action, _, err := cli.Client().Network.AddSubnet(cli.Context, network, opts)
	if err != nil {
		return err
	}
	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Subnet added to network %d\n", network.ID)

	return nil
}
