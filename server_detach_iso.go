package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerDetachISOCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "detach-iso [FLAGS] SERVER",
		Short:                 "Detaches an ISO from the server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerDetachISO),
	}

	return cmd
}

func runServerDetachISO(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	action, _, err := cli.Client().Server.DetachISO(cli.Context, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("ISO detached from server %s\n", idOrName)
	return nil
}
