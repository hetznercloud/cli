package network

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var AddSubnetCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "add-subnet [options] --type <cloud|server|vswitch> --network-zone <zone> <network>",
			Short:                 "Add a subnet to a network",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("type", "", "Type of subnet (required)")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("cloud", "server", "vswitch"))
		_ = cmd.MarkFlagRequired("type")

		cmd.Flags().String("network-zone", "", "Name of network zone (required)")
		_ = cmd.RegisterFlagCompletionFunc("network-zone", cmpl.SuggestCandidatesF(client.Location().NetworkZones))
		_ = cmd.MarkFlagRequired("network-zone")

		cmd.Flags().IPNet("ip-range", net.IPNet{}, "Range to allocate IPs from")

		cmd.Flags().Int64("vswitch-id", 0, "ID of the vSwitch")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		subnetType, _ := cmd.Flags().GetString("type")
		networkZone, _ := cmd.Flags().GetString("network-zone")
		ipRange, _ := cmd.Flags().GetIPNet("ip-range")
		vSwitchID, _ := cmd.Flags().GetInt64("vswitch-id")
		idOrName := args[0]

		network, _, err := s.Client().Network().Get(s, idOrName)
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
		action, _, err := s.Client().Network().AddSubnet(s, network, opts)
		if err != nil {
			return err
		}
		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}
		cmd.Printf("Subnet added to network %d\n", network.ID)

		return nil
	},
}
