package server

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "server",
		Short:                 "Manage servers",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(cli.Context, client, cli),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		DeleteCmd.CobraCommand(cli.Context, client, cli, cli),
		RebootCmd.CobraCommand(cli.Context, client, cli, cli),
		PoweronCmd.CobraCommand(cli.Context, client, cli, cli),
		PoweroffCmd.CobraCommand(cli.Context, client, cli, cli),
		ResetCmd.CobraCommand(cli.Context, client, cli, cli),
		ShutdownCmd.CobraCommand(cli.Context, client, cli, cli),
		CreateImageCmd.CobraCommand(cli.Context, client, cli, cli),
		ResetPasswordCmd.CobraCommand(cli.Context, client, cli, cli),
		EnableRescueCmd.CobraCommand(cli.Context, client, cli, cli),
		DisableRescueCmd.CobraCommand(cli.Context, client, cli, cli),
		AttachISOCmd.CobraCommand(cli.Context, client, cli, cli),
		DetachISOCmd.CobraCommand(cli.Context, client, cli, cli),
		UpdateCmd.CobraCommand(cli.Context, client, cli),
		ChangeTypeCmd.CobraCommand(cli.Context, client, cli, cli),
		RebuildCmd.CobraCommand(cli.Context, client, cli, cli),
		EnableBackupCmd.CobraCommand(cli.Context, client, cli, cli),
		DisableBackupCmd.CobraCommand(cli.Context, client, cli, cli),
		EnableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		SSHCmd.CobraCommand(cli.Context, client, cli, cli),
		LabelCmds.AddCobraCommand(cli.Context, client, cli),
		LabelCmds.RemoveCobraCommand(cli.Context, client, cli),
		SetRDNSCmd.CobraCommand(cli.Context, client, cli, cli),
		AttachToNetworkCmd.CobraCommand(cli.Context, client, cli, cli),
		DetachFromNetworkCmd.CobraCommand(cli.Context, client, cli, cli),
		ChangeAliasIPsCmd.CobraCommand(cli.Context, client, cli, cli),
		IPCmd.CobraCommand(cli.Context, client, cli, cli),
		RequestConsoleCmd.CobraCommand(cli.Context, client, cli, cli),
		MetricsCmd.CobraCommand(cli.Context, client, cli, cli),
		AddToPlacementGroupCmd.CobraCommand(cli.Context, client, cli, cli),
		RemoveFromPlacementGroupCmd.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}
