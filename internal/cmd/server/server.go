package server

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "server",
		Short:                 "Manage servers",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddGroup(&cobra.Group{ID: "general", Title: "General"})
	cmd.AddGroup(&cobra.Group{ID: "protection", Title: "Protection"})
	cmd.AddGroup(&cobra.Group{ID: "rescue", Title: "Rescue"})
	cmd.AddGroup(&cobra.Group{ID: "power", Title: "Power/Reboot"})
	cmd.AddGroup(&cobra.Group{ID: "network", Title: "Networks"})
	cmd.AddGroup(&cobra.Group{ID: "iso", Title: "ISO"})
	cmd.AddGroup(&cobra.Group{ID: "placement-group", Title: "Placement Groups"})
	cmd.AddGroup(&cobra.Group{ID: "backup", Title: "Backup"})

	cmd.AddCommand(
		util.WithGroup("general", ListCmd.CobraCommand(s)),
		util.WithGroup("general", DescribeCmd.CobraCommand(s)),
		util.WithGroup("general", CreateCmd.CobraCommand(s)),
		util.WithGroup("general", DeleteCmd.CobraCommand(s)),
		util.WithGroup("general", UpdateCmd.CobraCommand(s)),
		util.WithGroup("general", CreateImageCmd.CobraCommand(s)),
		util.WithGroup("general", ChangeTypeCmd.CobraCommand(s)),
		util.WithGroup("general", RebuildCmd.CobraCommand(s)),
		util.WithGroup("general", LabelCmds.AddCobraCommand(s)),
		util.WithGroup("general", LabelCmds.RemoveCobraCommand(s)),

		util.WithGroup("protection", EnableProtectionCmd.CobraCommand(s)),
		util.WithGroup("protection", DisableProtectionCmd.CobraCommand(s)),

		util.WithGroup("rescue", EnableRescueCmd.CobraCommand(s)),
		util.WithGroup("rescue", DisableRescueCmd.CobraCommand(s)),

		util.WithGroup("power", PoweronCmd.CobraCommand(s)),
		util.WithGroup("power", PoweroffCmd.CobraCommand(s)),
		util.WithGroup("power", RebootCmd.CobraCommand(s)),
		util.WithGroup("power", ShutdownCmd.CobraCommand(s)),
		util.WithGroup("power", ResetCmd.CobraCommand(s)),

		util.WithGroup("network", AttachToNetworkCmd.CobraCommand(s)),
		util.WithGroup("network", DetachFromNetworkCmd.CobraCommand(s)),
		util.WithGroup("network", ChangeAliasIPsCmd.CobraCommand(s)),
		util.WithGroup("network", SetRDNSCmd.CobraCommand(s)),

		util.WithGroup("iso", AttachISOCmd.CobraCommand(s)),
		util.WithGroup("iso", DetachISOCmd.CobraCommand(s)),

		util.WithGroup("placement-group", AddToPlacementGroupCmd.CobraCommand(s)),
		util.WithGroup("placement-group", RemoveFromPlacementGroupCmd.CobraCommand(s)),

		util.WithGroup("backup", EnableBackupCmd.CobraCommand(s)),
		util.WithGroup("backup", DisableBackupCmd.CobraCommand(s)),

		util.WithGroup("", SSHCmd.CobraCommand(s)),
		util.WithGroup("", IPCmd.CobraCommand(s)),
		util.WithGroup("", RequestConsoleCmd.CobraCommand(s)),
		util.WithGroup("", ResetPasswordCmd.CobraCommand(s)),
		util.WithGroup("", MetricsCmd.CobraCommand(s)),
	)
	return cmd
}
