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
		deleteCmd.CobraCommand(cli.Context, client, cli),
		newRebootCommand(cli),
		newPoweronCommand(cli),
		newPoweroffCommand(cli),
		newResetCommand(cli),
		newShutdownCommand(cli),
		newCreateImageCommand(cli),
		newResetPasswordCommand(cli),
		EnableRescueCommand.CobraCommand(cli.Context, client, cli, cli),
		newDisableRescueCommand(cli),
		AttachISOCommand.CobraCommand(cli.Context, client, cli, cli),
		newDetachISOCommand(cli),
		updateCmd.CobraCommand(cli.Context, client, cli),
		newChangeTypeCommand(cli),
		RebuildCommand.CobraCommand(cli.Context, client, cli, cli),
		newEnableBackupCommand(cli),
		newDisableBackupCommand(cli),
		newEnableProtectionCommand(cli),
		newDisableProtectionCommand(cli),
		newSSHCommand(cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
		setRDNSCmd.CobraCommand(cli.Context, client, cli, cli),
		AttachToNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		DetachFromNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		ChangeAliasIPsCommand.CobraCommand(cli.Context, client, cli, cli),
		newIPCommand(cli),
		newRequestConsoleCommand(cli),
		newMetricsCommand(cli),
		AddToPlacementGroupCommand.CobraCommand(cli.Context, client, cli, cli),
		RemoveFromPlacementGroup.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}
