package floatingip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "floating-ip",
		Short:                 "Manage Floating IPs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		UpdateCmd.CobraCommand(cli.Context, client, cli),
		ListCmd.CobraCommand(cli.Context, client, cli),
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		AssignCmd.CobraCommand(cli.Context, client, cli, cli),
		UnassignCmd.CobraCommand(cli.Context, client, cli, cli),
		DeleteCmd.CobraCommand(cli.Context, client, cli, cli),
		EnableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		LabelCmds.AddCobraCommand(cli.Context, client, cli),
		LabelCmds.RemoveCobraCommand(cli.Context, client, cli),
		SetRDNSCmd.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}
