package server

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newResetPasswordCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "reset-password [FLAGS] SERVER",
		Short:                 "Reset the root password of a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runResetPassword),
	}
	return cmd
}

func runResetPassword(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	result, _, err := cli.Client().Server.ResetPassword(cli.Context, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}

	fmt.Printf("Password of server %d reset to: %s\n", server.ID, result.RootPassword)
	return nil
}
