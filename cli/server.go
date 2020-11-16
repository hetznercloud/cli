package cli

import "github.com/spf13/cobra"

func newServerCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "server",
		Short:                 "Manage servers",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newServerListCommand(cli),
		newServerDescribeCommand(cli),
		newServerCreateCommand(cli),
		newServerDeleteCommand(cli),
		newServerRebootCommand(cli),
		newServerPoweronCommand(cli),
		newServerPoweroffCommand(cli),
		newServerResetCommand(cli),
		newServerShutdownCommand(cli),
		newServerCreateImageCommand(cli),
		newServerResetPasswordCommand(cli),
		newServerEnableRescueCommand(cli),
		newServerDisableRescueCommand(cli),
		newServerAttachISOCommand(cli),
		newServerDetachISOCommand(cli),
		newServerUpdateCommand(cli),
		newServerChangeTypeCommand(cli),
		newServerRebuildCommand(cli),
		newServerEnableBackupCommand(cli),
		newServerDisableBackupCommand(cli),
		newServerEnableProtectionCommand(cli),
		newServerDisableProtectionCommand(cli),
		newServerSSHCommand(cli),
		newServerAddLabelCommand(cli),
		newServerRemoveLabelCommand(cli),
		newServerSetRDNSCommand(cli),
		newServerAttachToNetworkCommand(cli),
		newServerDetachFromNetworkCommand(cli),
		newServerChangeAliasIPsCommand(cli),
		newServerIPCommand(cli),
		newServerRequestConsoleCommand(cli),
		newServerMetricsCommand(cli),
	)
	return cmd
}
