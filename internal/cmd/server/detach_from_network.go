package server

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"

	"github.com/spf13/cobra"
)

func newDetachFromNetworkCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "detach-from-network [FLAGS] SERVER",
		Short:                 "Detach a server from a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerDetachFromNetwork),
	}
	cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
	cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(cli.NetworkNames))
	cmd.MarkFlagRequired("network")

	return cmd
}

func runServerDetachFromNetwork(cli *state.State, cmd *cobra.Command, args []string) error {
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

	opts := hcloud.ServerDetachFromNetworkOpts{
		Network: network,
	}
	action, _, err := cli.Client().Server.DetachFromNetwork(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Server %d detached from network %d\n", server.ID, network.ID)
	return nil
}
