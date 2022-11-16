package placementgroup

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "placement-group",
		Short:                 "Manage Placement Groups",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		ListCmd.CobraCommand(cli.Context, client, cli, cli.Config.SubcommandDefaults[cmd.Use]),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		UpdateCmd.CobraCommand(cli.Context, client, cli),
		DeleteCmd.CobraCommand(cli.Context, client, cli),
		LabelCmds.AddCobraCommand(cli.Context, client, cli),
		LabelCmds.RemoveCobraCommand(cli.Context, client, cli),
	)
	return cmd
}
