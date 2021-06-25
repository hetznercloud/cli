package firewall

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "firewall",
		Short:                 "Manage Firewalls",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		listCmd.CobraCommand(cli.Context, client, cli),
		describeCmd.CobraCommand(cli.Context, client, cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		newReplaceRulesCommand(cli),
		deleteCmd.CobraCommand(cli.Context, client, cli),
		newAddRuleCommand(cli),
		newDeleteRuleCommand(cli),
		newApplyToResourceCommand(cli),
		newRemoveFromResourceCommand(cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
	)
	return cmd
}
