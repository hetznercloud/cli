package primaryip

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "primary-ip",
		Short:                 "Manage Primary IPs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(cli.Context, client, cli),
		describeCmd.CobraCommand(cli.Context, client, cli),
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		updateCmd.CobraCommand(cli.Context, client, cli),
		deleteCmd.CobraCommand(cli.Context, client, cli, cli),
		AssignCmd.CobraCommand(cli.Context, client, cli, cli),
		UnAssignCmd.CobraCommand(cli.Context, client, cli, cli),
		ChangeDNSCmd.CobraCommand(cli.Context, client, cli, cli),
		EnableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
	)
	return cmd
}
