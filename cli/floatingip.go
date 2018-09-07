package cli

import "github.com/spf13/cobra"

func newFloatingIPCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "floating-ip",
		Short:                 "Manage Floating IPs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runFloatingIP),
	}
	cmd.AddCommand(
		newFloatingIPUpdateCommand(cli),
		newFloatingIPListCommand(cli),
		newFloatingIPCreateCommand(cli),
		newFloatingIPDescribeCommand(cli),
		newFloatingIPAssignCommand(cli),
		newFloatingIPUnassignCommand(cli),
		newFloatingIPDeleteCommand(cli),
		newFloatingIPEnableProtectionCommand(cli),
		newFloatingIPDisableProtectionCommand(cli),
		newFloatingIPAddLabelCommand(cli),
		newFloatingIPRemoveLabelCommand(cli),
	)
	return cmd
}

func runFloatingIP(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
