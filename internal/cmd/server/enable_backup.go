package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var EnableBackupCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "enable-backup <server>",
			Short:                 "Enable backup for a server",
			Args:                  util.ValidateExact,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String(
			"window", "",
			"(deprecated) The time window for the daily backup to run. All times are in UTC. 22-02 means that the backup will be started between 10 PM and 2 AM.")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		window, _ := cmd.Flags().GetString("window")
		if window != "" {
			cmd.Print("[WARN] The ability to specify a backup window when enabling backups has been removed. Ignoring flag.\n")
		}

		action, _, err := s.Client().Server().EnableBackup(s, server, "")
		if err != nil {
			return err
		}

		if err := s.ActionProgress(cmd, s, action); err != nil {
			return err
		}

		cmd.Printf("Backup enabled for server %d\n", server.ID)
		return nil
	},
}
