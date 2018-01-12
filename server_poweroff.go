package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerPoweroffCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "poweroff [FLAGS] SERVER",
		Short:                 "Poweroff a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerPoweroff),
	}
	return cmd
}

func runServerPoweroff(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	action, _, err := cli.Client().Server.Poweroff(cli.Context, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Server %s stopped\n", server.Name)
	return nil
}
