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

var ChangeIPRangeCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "change-ip-range --ip-range <ip-range> <network>",
			Short:                 "Change the IP range of a network",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().IPNet("ip-range", net.IPNet{}, "New IP range (required)")
		cmd.MarkFlagRequired("ip-range")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		network, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", idOrName)
		}

		ipRange, _ := cmd.Flags().GetIPNet("ip-range")
		opts := hcloud.NetworkChangeIPRangeOpts{
			IPRange: &ipRange,
		}

		action, _, err := s.Client().Network().ChangeIPRange(s, network, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}
		cmd.Printf("IP range of network %d changed\n", network.ID)
		return nil
	},
}
