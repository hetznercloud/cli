package firewall

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "firewall",
		Short:                 "Manage Firewalls",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newListCommand(cli),
		newDescribeCommand(cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		newReplaceRulesCommand(cli),
		newAddLabelCommand(cli),
		newRemoveLabelCommand(cli),
		newDeleteCommand(cli),
		newAddRuleCommand(cli),
		newDeleteRuleCommand(cli),
		newApplyToResourceCommand(cli),
		newRemoveFromResourceCommand(cli),
	)
	return cmd
}
