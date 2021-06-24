package floatingip

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
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
		newUpdateCommand(cli),
		newListCommand(cli),
		newCreateCommand(cli),
		newDescribeCommand(cli),
		newAssignCommand(cli),
		newUnassignCommand(cli),
		newDeleteCommand(cli),
		newEnableProtectionCommand(cli),
		newDisableProtectionCommand(cli),
		newAddLabelCommand(cli),
		newRemoveLabelCommand(cli),
		newSetRDNSCommand(cli),
	)
	return cmd
}
