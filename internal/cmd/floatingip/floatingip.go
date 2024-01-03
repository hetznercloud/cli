package floatingip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "floating-ip",
		Short:                 "Manage Floating IPs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		UpdateCmd.CobraCommand(s),
		ListCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		AssignCmd.CobraCommand(s),
		UnassignCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		SetRDNSCmd.CobraCommand(s),
	)
	return cmd
}
