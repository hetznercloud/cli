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

var ChangeAliasIPsCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "change-alias-ips [options] --network <network> <server>",
			Short:                 "Change a server's alias IPs in a network",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
		_ = cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(client.Network().Names))
		_ = cmd.MarkFlagRequired("network")

		cmd.Flags().StringSlice("alias-ips", nil, "New alias IPs")
		cmd.Flags().Bool("clear", false, "Remove all alias IPs")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		clear, _ := cmd.Flags().GetBool("clear")
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		networkIDOrName, _ := cmd.Flags().GetString("network")
		network, _, err := s.Client().Network().Get(s, networkIDOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", networkIDOrName)
		}

		aliasIPs, _ := cmd.Flags().GetStringSlice("alias-ips")

		opts := hcloud.ServerChangeAliasIPsOpts{
			Network: network,
		}
		if clear {
			opts.AliasIPs = []net.IP{}
		} else {
			for _, aliasIP := range aliasIPs {
				opts.AliasIPs = append(opts.AliasIPs, net.ParseIP(aliasIP))
			}
		}
		action, _, err := s.Client().Server().ChangeAliasIPs(s, server, opts)

		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Alias IPs changed for server %d in network %d\n", server.ID, network.ID)
		return nil
	},
}
