package server

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

var AttachToNetworkCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "attach-to-network [options] --network <network> <server>",
			Short:                 "Attach a Server to a Network",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
		_ = cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(client.Network().Names))
		_ = cmd.MarkFlagRequired("network")

		cmd.Flags().IP("ip", nil, "IP address to assign to the Server (auto-assigned if omitted)")
		cmd.Flags().IPNet("ip-range", net.IPNet{}, "IP range in CIDR block notation of the subnet to attach to")
		cmd.Flags().IPSlice("alias-ips", []net.IP{}, "Additional IP addresses to be assigned to the Server")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}

		networkIDOrName, _ := cmd.Flags().GetString("network")
		network, _, err := s.Client().Network().Get(s, networkIDOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("Network not found: %s", networkIDOrName)
		}

		ip, _ := cmd.Flags().GetIP("ip")
		ipRange, _ := cmd.Flags().GetIPNet("ip-range")
		aliasIPs, _ := cmd.Flags().GetIPSlice("alias-ips")

		opts := hcloud.ServerAttachToNetworkOpts{
			Network: network,
			IP:      ip,
		}
		opts.AliasIPs = append(opts.AliasIPs, aliasIPs...)
		if cmd.Flags().Changed("ip-range") {
			opts.IPRange = &ipRange
		}
		action, _, err := s.Client().Server().AttachToNetwork(s, server, opts)

		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Server %d attached to Network %d\n", server.ID, network.ID)
		return nil
	},
}
