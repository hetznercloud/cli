package server

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newIPCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ip SERVER FLAGS",
		Short:                 "Print a server's IP address",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runIP),
	}
	cmd.Flags().BoolP("ipv6", "6", false, "Print the first address of the IPv6 public server network")
	return cmd
}

func runIP(cli *state.State, cmd *cobra.Command, args []string) error {
	ipv6, err := cmd.Flags().GetBool("ipv6")
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}
	if ipv6 {
		fmt.Println(server.PublicNet.IPv6.IP.String() + "1")
	} else {
		fmt.Println(server.PublicNet.IPv4.IP.String())
	}
	return nil
}
