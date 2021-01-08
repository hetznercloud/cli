package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newServerEnableBackupCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "enable-backup [FLAGS] SERVER",
		Short:                 "Enable backup for a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerEnableBackup),
	}
	cmd.Flags().String(
		"window", "",
		"(deprecated) The time window for the daily backup to run. All times are in UTC. 22-02 means that the backup will be started between 10 PM and 2 AM.")
	return cmd
}

func runServerEnableBackup(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	window, _ := cmd.Flags().GetString("window")
	if window != "" {
		fmt.Print("[WARN] The ability to specify a backup window when enabling backups has been removed. Ignoring flag.\n")
	}

	action, _, err := cli.Client().Server.EnableBackup(cli.Context, server, "")
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Backup enabled for server %d\n", server.ID)
	return nil
}
