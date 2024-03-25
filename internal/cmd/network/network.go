package network

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "network",
		Short:                 "Manage networks",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	util.AddGroup(cmd, "general", "General",
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		ChangeIPRangeCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "route", "Routes",
		AddRouteCmd.CobraCommand(s),
		RemoveRouteCmd.CobraCommand(s),
		ExposeRoutesToVSwitchCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "subnet", "Subnets",
		AddSubnetCmd.CobraCommand(s),
		RemoveSubnetCmd.CobraCommand(s),
	)
	return cmd
}
