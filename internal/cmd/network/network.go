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
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddGroup(&cobra.Group{ID: "general", Title: "General"})
	cmd.AddGroup(&cobra.Group{ID: "protection", Title: "Protection"})
	cmd.AddGroup(&cobra.Group{ID: "route", Title: "Routes"})
	cmd.AddGroup(&cobra.Group{ID: "subnet", Title: "Subnets"})

	cmd.AddCommand(
		util.WithGroup("general", ListCmd.CobraCommand(s)),
		util.WithGroup("general", DescribeCmd.CobraCommand(s)),
		util.WithGroup("general", DeleteCmd.CobraCommand(s)),
		util.WithGroup("general", CreateCmd.CobraCommand(s)),
		util.WithGroup("general", UpdateCmd.CobraCommand(s)),
		util.WithGroup("general", LabelCmds.AddCobraCommand(s)),
		util.WithGroup("general", LabelCmds.RemoveCobraCommand(s)),
		util.WithGroup("general", ChangeIPRangeCmd.CobraCommand(s)),

		util.WithGroup("protection", EnableProtectionCmd.CobraCommand(s)),
		util.WithGroup("protection", DisableProtectionCmd.CobraCommand(s)),

		util.WithGroup("route", AddRouteCmd.CobraCommand(s)),
		util.WithGroup("route", RemoveRouteCmd.CobraCommand(s)),
		util.WithGroup("route", ExposeRoutesToVSwitchCmd.CobraCommand(s)),

		util.WithGroup("subnet", AddSubnetCmd.CobraCommand(s)),
		util.WithGroup("subnet", RemoveSubnetCmd.CobraCommand(s)),
	)
	return cmd
}
