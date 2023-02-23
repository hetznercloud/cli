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
		updateCmd.CobraCommand(cli.Context, client, cli),
		ReplaceRulesCommand.CobraCommand(cli.Context, client, cli, cli),
		deleteCmd.CobraCommand(cli.Context, client, cli),
		AddRuleCommand.CobraCommand(cli.Context, client, cli, cli),
		DeleteRuleCommand.CobraCommand(cli.Context, client, cli, cli),
		ApplyToResourceCommand.CobraCommand(cli.Context, client, cli, cli),
		RemoveFromResourceCommand.CobraCommand(cli.Context, client, cli, cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
	)
	return cmd
}
