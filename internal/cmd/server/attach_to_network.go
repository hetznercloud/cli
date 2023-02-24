package server

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

var AttachToNetworkCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "attach-to-network [FLAGS] SERVER",
			Short:                 "Attach a server to a network",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
		cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(client.Network().Names))
		cmd.MarkFlagRequired("network")

		cmd.Flags().IP("ip", nil, "IP address to assign to the server (auto-assigned if omitted)")
		cmd.Flags().IPSlice("alias-ips", []net.IP{}, "Additional IP addresses to be assigned to the server")

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		networkIDOrName, _ := cmd.Flags().GetString("network")
		network, _, err := client.Network().Get(ctx, networkIDOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", networkIDOrName)
		}

		ip, _ := cmd.Flags().GetIP("ip")
		aliasIPs, _ := cmd.Flags().GetIPSlice("alias-ips")

		opts := hcloud.ServerAttachToNetworkOpts{
			Network: network,
			IP:      ip,
		}
		for _, aliasIP := range aliasIPs {
			opts.AliasIPs = append(opts.AliasIPs, aliasIP)
		}
		action, _, err := client.Server().AttachToNetwork(ctx, server, opts)

		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Server %d attached to network %d\n", server.ID, network.ID)
		return nil
	},
}
