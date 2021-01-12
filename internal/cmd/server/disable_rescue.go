package server

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newDisableRescueCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disable-rescue [FLAGS] SERVER",
		Short:                 "Disable rescue for a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDisableRescue),
	}
	return cmd
}

func runDisableRescue(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	action, _, err := cli.Client().Server.DisableRescue(cli.Context, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Rescue disabled for server %d\n", server.ID)
	return nil
}
