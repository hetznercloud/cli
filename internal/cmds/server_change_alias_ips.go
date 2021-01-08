package cmds

import (
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerChangeAliasIPsCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "change-alias-ips [FLAGS] SERVER",
		Short:                 "Change a server's alias IPs in a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerChangeAliasIPsk),
	}

	cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
	cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(cli.NetworkNames))
	cmd.MarkFlagRequired("network")

	cmd.Flags().StringSlice("alias-ips", nil, "New alias IPs")
	cmd.Flags().Bool("clear", false, "Remove all alias IPs")

	return cmd
}

func runServerChangeAliasIPsk(cli *state.State, cmd *cobra.Command, args []string) error {
	clear, _ := cmd.Flags().GetBool("clear")
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	networkIDOrName, _ := cmd.Flags().GetString("network")
	network, _, err := cli.Client().Network.Get(cli.Context, networkIDOrName)
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
	action, _, err := cli.Client().Server.ChangeAliasIPs(cli.Context, server, opts)

	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Alias IPs changed for server %d in network %d\n", server.ID, network.ID)
	return nil
}
