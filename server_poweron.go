package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerPoweronCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "poweron [FLAGS] SERVER",
		Short:                 "Poweron a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerPoweron),
	}
	return cmd
}

func runServerPoweron(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	action, _, err := cli.Client().Server.Poweron(cli.Context, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Server %s started\n", server.Name)
	return nil
}
