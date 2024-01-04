package network

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "network",
		Short:                 "Manage networks",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		ChangeIPRangeCmd.CobraCommand(s),
		AddRouteCmd.CobraCommand(s),
		RemoveRouteCmd.CobraCommand(s),
		AddSubnetCmd.CobraCommand(s),
		RemoveSubnetCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
		ExposeRoutesToVSwitchCmd.CobraCommand(s),
	)
	return cmd
}
