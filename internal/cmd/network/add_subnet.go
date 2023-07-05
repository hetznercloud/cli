package network

import (
	"context"
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var AddSubnetCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "add-subnet NETWORK FLAGS",
			Short:                 "Add a subnet to a network",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("type", "", "Type of subnet (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("cloud", "server", "vswitch"))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("network-zone", "", "Name of network zone (required)")
		cmd.RegisterFlagCompletionFunc("network-zone", cmpl.SuggestCandidatesF(client.Location().NetworkZones))
		cmd.MarkFlagRequired("network-zone")

		cmd.Flags().IPNet("ip-range", net.IPNet{}, "Range to allocate IPs from")

		cmd.Flags().Int64("vswitch-id", 0, "ID of the vSwitch")
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		subnetType, _ := cmd.Flags().GetString("type")
		networkZone, _ := cmd.Flags().GetString("network-zone")
		ipRange, _ := cmd.Flags().GetIPNet("ip-range")
		vSwitchID, _ := cmd.Flags().GetInt64("vswitch-id")
		idOrName := args[0]

		network, _, err := client.Network().Get(ctx, idOrName)
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
		if subnetType == "vswitch" {
			subnet.VSwitchID = vSwitchID
		}

		opts := hcloud.NetworkAddSubnetOpts{
			Subnet: subnet,
		}
		action, _, err := client.Network().AddSubnet(ctx, network, opts)
		if err != nil {
			return err
		}
		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}
		fmt.Printf("Subnet added to network %d\n", network.ID)

		return nil
	},
}
