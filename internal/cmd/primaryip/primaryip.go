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

	util.AddGroup(cmd, "general", "General",
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "assign", "Assign",
		AssignCmd.CobraCommand(s),
		UnAssignCmd.CobraCommand(s),
	)

	cmd.AddCommand(SetRDNSCmd.CobraCommand(s))
	return cmd
}
