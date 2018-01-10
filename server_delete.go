package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] SERVER",
		Short:                 "Delete a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerDelete),
	}
	return cmd
}

func runServerDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	_, err = cli.Client().Server.Delete(cli.Context, server)
	if err != nil {
		return err
	}

	fmt.Printf("Server %s deleted\n", idOrName)
	return nil
}
