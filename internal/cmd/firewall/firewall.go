package firewall

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
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
		ListCmd.CobraCommand(cli.Context, client, cli),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		UpdateCmd.CobraCommand(cli.Context, client, cli),
		ReplaceRulesCmd.CobraCommand(cli.Context, client, cli, cli),
		DeleteCmd.CobraCommand(cli.Context, client, cli, cli),
		AddRuleCmd.CobraCommand(cli.Context, client, cli, cli),
		DeleteRuleCmd.CobraCommand(cli.Context, client, cli, cli),
		ApplyToResourceCmd.CobraCommand(cli.Context, client, cli, cli),
		RemoveFromResourceCmd.CobraCommand(cli.Context, client, cli, cli),
		LabelCmds.AddCobraCommand(cli.Context, client, cli),
		LabelCmds.RemoveCobraCommand(cli.Context, client, cli),
	)
	return cmd
}
