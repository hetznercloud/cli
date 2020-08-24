package cli

import "github.com/spf13/cobra"

func newFloatingIPCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "floating-ip",
		Short:                 "Manage Floating IPs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
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
		newFloatingIPSetRDNSCommand(cli),
	)
	return cmd
}
