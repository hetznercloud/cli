package network

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "network",
		Short:                 "Manage networks",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		newDescribeCommand(cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		newDeleteCommand(cli),
		newChangeIPRangeCommand(cli),
		newAddRouteCommand(cli),
		newRemoveRouteCommand(cli),
		newAddSubnetCommand(cli),
		newRemoveSubnetCommand(cli),
		newAddLabelCommand(cli),
		newRemoveLabelCommand(cli),
		newEnableProtectionCommand(cli),
		newDisableProtectionCommand(cli),
	)
	return cmd
}
