package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerDisableBackupCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disable-backup [FLAGS] SERVER",
		Short:                 "Disable backups of a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerDisableBackup),
	}
	return cmd
}

func runServerDisableBackup(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	action, _, err := cli.Client().Server.DisableBackup(cli.Context, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Backups of server %s disabled\n", idOrName)
	return nil
}
