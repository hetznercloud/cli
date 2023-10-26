package network

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "network",
		Short:                 "Manage networks",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(cli.Context, client, cli),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		UpdateCmd.CobraCommand(cli.Context, client, cli),
		DeleteCmd.CobraCommand(cli.Context, client, cli, cli),
		ChangeIPRangeCmd.CobraCommand(cli.Context, client, cli, cli),
		AddRouteCmd.CobraCommand(cli.Context, client, cli, cli),
		RemoveRouteCmd.CobraCommand(cli.Context, client, cli, cli),
		AddSubnetCmd.CobraCommand(cli.Context, client, cli, cli),
		RemoveSubnetCmd.CobraCommand(cli.Context, client, cli, cli),
		LabelCmds.AddCobraCommand(cli.Context, client, cli),
		LabelCmds.RemoveCobraCommand(cli.Context, client, cli),
		EnableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		ExposeRoutesToVSwitchCmd.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}
