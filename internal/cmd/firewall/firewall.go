package firewall

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
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
		describeCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		newReplaceRulesCommand(cli),
		deleteCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		newAddRuleCommand(cli),
		newDeleteRuleCommand(cli),
		newApplyToResourceCommand(cli),
		newRemoveFromResourceCommand(cli),
		labelCmds.AddCobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		labelCmds.RemoveCobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
	)
	return cmd
}
