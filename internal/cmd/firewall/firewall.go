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
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddGroup(&cobra.Group{ID: "general", Title: "General"})
	cmd.AddGroup(&cobra.Group{ID: "rule", Title: "Rules"})
	cmd.AddGroup(&cobra.Group{ID: "resource", Title: "Resources"})

	cmd.AddCommand(
		util.WithGroup("general", ListCmd.CobraCommand(s)),
		util.WithGroup("general", DescribeCmd.CobraCommand(s)),
		util.WithGroup("general", CreateCmd.CobraCommand(s)),
		util.WithGroup("general", DeleteCmd.CobraCommand(s)),
		util.WithGroup("general", UpdateCmd.CobraCommand(s)),
		util.WithGroup("general", LabelCmds.AddCobraCommand(s)),
		util.WithGroup("general", LabelCmds.RemoveCobraCommand(s)),

		util.WithGroup("rule", ReplaceRulesCmd.CobraCommand(s)),
		util.WithGroup("rule", AddRuleCmd.CobraCommand(s)),
		util.WithGroup("rule", DeleteRuleCmd.CobraCommand(s)),

		util.WithGroup("resource", ApplyToResourceCmd.CobraCommand(s)),
		util.WithGroup("resource", RemoveFromResourceCmd.CobraCommand(s)),
	)
	return cmd
}
