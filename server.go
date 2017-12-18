package cli

import "github.com/spf13/cobra"

func newServerCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "server",
		Short:            "Manage servers",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE:             cli.wrap(runServer),
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
		newServerCreateimageCommand(cli),
		newServerResetpasswordCommand(cli),
		newServerEnableRescueCommand(cli),
		newServerDisableRescueCommand(cli),
	)
	return cmd
}

func runServer(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
