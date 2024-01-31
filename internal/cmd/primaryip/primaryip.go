package primaryip

import (
	"github.com/spf13/cobra"

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
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		AssignCmd.CobraCommand(s),
		UnAssignCmd.CobraCommand(s),
		SetRDNSCmd.CobraCommand(s),
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)
	return cmd
}
