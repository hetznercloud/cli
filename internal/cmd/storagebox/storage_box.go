package storagebox

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/snapshot"
	"github.com/hetznercloud/cli/internal/cmd/storagebox/subaccount"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "storage-box",
		Aliases:               []string{"storage-boxes"},
		Short:                 "Manage Storage Boxes",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	util.AddGroup(cmd, "general", "General",
		ListCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		FoldersCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		UpdateCmd.CobraCommand(s),
		ChangeTypeCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "snapshots", "Snapshots",
		snapshot.NewCommand(s),
		EnableSnapshotPlanCmd.CobraCommand(s),
		DisableSnapshotPlanCmd.CobraCommand(s),
		RollbackSnapshotCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "account", "Account",
		subaccount.NewCommand(s),
		ResetPasswordCmd.CobraCommand(s),
		UpdateAccessSettingsCmd.CobraCommand(s),
	)
	return cmd
}
