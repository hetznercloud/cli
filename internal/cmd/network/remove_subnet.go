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

var RemoveSubnetCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "remove-subnet --ip-range <ip-range> <network>",
			Short:                 "Remove a subnet from a Network",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().IPNet("ip-range", net.IPNet{}, "Subnet IP range (required)")
		_ = cmd.MarkFlagRequired("ip-range")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		ipRange, _ := cmd.Flags().GetIPNet("ip-range")
		idOrName := args[0]
		network, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("Network not found: %s", idOrName)
		}

		opts := hcloud.NetworkDeleteSubnetOpts{
			Subnet: hcloud.NetworkSubnet{
				IPRange: &ipRange,
			},
		}
		action, _, err := s.Client().Network().DeleteSubnet(s, network, opts)
		if err != nil {
			return err
		}
		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}
		cmd.Printf("Subnet %s removed from Network %d\n", ipRange.String(), network.ID)

		return nil
	},
}
