package primaryip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "primary-ip",
		Short:                 "Manage Primary IPs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddGroup(&cobra.Group{ID: "general", Title: "General"})
	cmd.AddGroup(&cobra.Group{ID: "protection", Title: "Protection"})
	cmd.AddGroup(&cobra.Group{ID: "assign", Title: "Assign"})

	cmd.AddCommand(
		util.WithGroup("general", ListCmd.CobraCommand(s)),
		util.WithGroup("general", DescribeCmd.CobraCommand(s)),
		util.WithGroup("general", DeleteCmd.CobraCommand(s)),
		util.WithGroup("general", CreateCmd.CobraCommand(s)),
		util.WithGroup("general", UpdateCmd.CobraCommand(s)),
		util.WithGroup("general", LabelCmds.AddCobraCommand(s)),
		util.WithGroup("general", LabelCmds.RemoveCobraCommand(s)),

		util.WithGroup("protection", EnableProtectionCmd.CobraCommand(s)),
		util.WithGroup("protection", DisableProtectionCmd.CobraCommand(s)),

		util.WithGroup("assign", AssignCmd.CobraCommand(s)),
		util.WithGroup("assign", UnAssignCmd.CobraCommand(s)),

		util.WithGroup("", SetRDNSCmd.CobraCommand(s)),
	)
	return cmd
}
