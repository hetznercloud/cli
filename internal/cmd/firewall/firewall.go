package firewall

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "firewall",
		Short:                 "Manage Firewalls",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		ReplaceRulesCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		AddRuleCmd.CobraCommand(s),
		DeleteRuleCmd.CobraCommand(s),
		ApplyToResourceCmd.CobraCommand(s),
		RemoveFromResourceCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)
	return cmd
}
