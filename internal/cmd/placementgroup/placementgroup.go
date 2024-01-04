package placementgroup

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "placement-group",
		Short:                 "Manage Placement Groups",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		CreateCmd.CobraCommand(s),
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)
	return cmd
}
