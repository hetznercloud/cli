package firewall

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "firewall",
		Short:                 "Manage Firewalls",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	util.AddGroup(cmd, "general", "General",
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)

	util.AddGroup(cmd, "rule", "Rules",
		ReplaceRulesCmd.CobraCommand(s),
		AddRuleCmd.CobraCommand(s),
		DeleteRuleCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "resource", "Resources",
		ApplyToResourceCmd.CobraCommand(s),
		RemoveFromResourceCmd.CobraCommand(s),
	)
	return cmd
}
