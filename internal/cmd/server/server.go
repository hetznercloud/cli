package server

import (
	"github.com/spf13/cobra"

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
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		RebootCmd.CobraCommand(s),
		PoweronCmd.CobraCommand(s),
		PoweroffCmd.CobraCommand(s),
		ResetCmd.CobraCommand(s),
		ShutdownCmd.CobraCommand(s),
		CreateImageCmd.CobraCommand(s),
		ResetPasswordCmd.CobraCommand(s),
		EnableRescueCmd.CobraCommand(s),
		DisableRescueCmd.CobraCommand(s),
		AttachISOCmd.CobraCommand(s),
		DetachISOCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		ChangeTypeCmd.CobraCommand(s),
		RebuildCmd.CobraCommand(s),
		EnableBackupCmd.CobraCommand(s),
		DisableBackupCmd.CobraCommand(s),
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
		SSHCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		SetRDNSCmd.CobraCommand(s),
		AttachToNetworkCmd.CobraCommand(s),
		DetachFromNetworkCmd.CobraCommand(s),
		ChangeAliasIPsCmd.CobraCommand(s),
		IPCmd.CobraCommand(s),
		RequestConsoleCmd.CobraCommand(s),
		MetricsCmd.CobraCommand(s),
		AddToPlacementGroupCmd.CobraCommand(s),
		RemoveFromPlacementGroupCmd.CobraCommand(s),
	)
	return cmd
}
