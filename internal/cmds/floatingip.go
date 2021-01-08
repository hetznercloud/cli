package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewFloatingIPCommand(cli *state.State) *cobra.Command {
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
