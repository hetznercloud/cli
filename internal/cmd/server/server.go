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

	util.AddGroup(cmd, "general", "General",
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		CreateImageCmd.CobraCommand(s),
		ChangeTypeCmd.CobraCommand(s),
		RebuildCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "rescue", "Rescue",
		EnableRescueCmd.CobraCommand(s),
		DisableRescueCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "power", "Power/Reboot",
		PoweronCmd.CobraCommand(s),
		PoweroffCmd.CobraCommand(s),
		RebootCmd.CobraCommand(s),
		ShutdownCmd.CobraCommand(s),
		ResetCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "network", "Networks",
		AttachToNetworkCmd.CobraCommand(s),
		DetachFromNetworkCmd.CobraCommand(s),
		ChangeAliasIPsCmd.CobraCommand(s),
		SetRDNSCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "iso", "ISO",
		AttachISOCmd.CobraCommand(s),
		DetachISOCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "placement-group", "Placement Groups",
		AddToPlacementGroupCmd.CobraCommand(s),
		RemoveFromPlacementGroupCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "backup", "Backup",
		EnableBackupCmd.CobraCommand(s),
		DisableBackupCmd.CobraCommand(s),
	)

	cmd.AddCommand(
		SSHCmd.CobraCommand(s),
		IPCmd.CobraCommand(s),
		RequestConsoleCmd.CobraCommand(s),
		ResetPasswordCmd.CobraCommand(s),
		MetricsCmd.CobraCommand(s),
	)
	return cmd
}
