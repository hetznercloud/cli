package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerResetCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "reset [FLAGS] SERVER",
		Short:                 "Reset a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerReset),
	}
	return cmd
}

func runServerReset(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	action, _, err := cli.Client().Server.Reset(cli.Context, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Server %s reset\n", server.Name)
	return nil
}
