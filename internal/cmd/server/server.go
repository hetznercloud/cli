package server

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
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
		describeCmd.CobraCommand(cli.Context, client, cli),
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		deleteCmd.CobraCommand(cli.Context, client, cli, cli),
		RebootCommand.CobraCommand(cli.Context, client, cli, cli),
		PoweronCommand.CobraCommand(cli.Context, client, cli, cli),
		PoweroffCommand.CobraCommand(cli.Context, client, cli, cli),
		ResetCommand.CobraCommand(cli.Context, client, cli, cli),
		ShutdownCommand.CobraCommand(cli.Context, client, cli, cli),
		newCreateImageCommand(cli),
		ResetPasswordCommand.CobraCommand(cli.Context, client, cli, cli),
		EnableRescueCommand.CobraCommand(cli.Context, client, cli, cli),
		DisableRescueCommand.CobraCommand(cli.Context, client, cli, cli),
		AttachISOCommand.CobraCommand(cli.Context, client, cli, cli),
		DetachISOCommand.CobraCommand(cli.Context, client, cli, cli),
		updateCmd.CobraCommand(cli.Context, client, cli),
		ChangeTypeCommand.CobraCommand(cli.Context, client, cli, cli),
		RebuildCommand.CobraCommand(cli.Context, client, cli, cli),
		EnableBackupCommand.CobraCommand(cli.Context, client, cli, cli),
		DisableBackupCommand.CobraCommand(cli.Context, client, cli, cli),
		EnableProtectionCommand.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCommand.CobraCommand(cli.Context, client, cli, cli),
		SSHCommand.CobraCommand(cli.Context, client, cli, cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
		setRDNSCmd.CobraCommand(cli.Context, client, cli, cli),
		AttachToNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		DetachFromNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		ChangeAliasIPsCommand.CobraCommand(cli.Context, client, cli, cli),
		IPCommand.CobraCommand(cli.Context, client, cli, cli),
		RequestConsoleCommand.CobraCommand(cli.Context, client, cli, cli),
		MetricsCommand.CobraCommand(cli.Context, client, cli, cli),
		AddToPlacementGroupCommand.CobraCommand(cli.Context, client, cli, cli),
		RemoveFromPlacementGroup.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}
