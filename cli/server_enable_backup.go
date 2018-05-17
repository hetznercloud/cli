package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newServerEnableBackupCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "enable-backup [FLAGS] SERVER",
		Short:                 "Enable backup for a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerEnableBackup),
	}
	cmd.Flags().String(
		"window", "",
		"The time window for the daily backup to run. All times are in UTC. 22-02 means that the backup will be started between 10 PM and 2 AM.")
	cmd.Flag("window").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_backup_windows"},
	}
	return cmd
}

func runServerEnableBackup(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	window, _ := cmd.Flags().GetString("window")
	action, _, err := cli.Client().Server.EnableBackup(cli.Context, server, window)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Backup enabled for server %s\n", server.Name)
	return nil
}
