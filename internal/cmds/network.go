package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewNetworkCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "network",
		Short:                 "Manage networks",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newNetworkListCommand(cli),
		newNetworkDescribeCommand(cli),
		newNetworkCreateCommand(cli),
		newNetworkUpdateCommand(cli),
		newNetworkDeleteCommand(cli),
		newNetworkChangeIPRangeCommand(cli),
		newNetworkAddRouteCommand(cli),
		newNetworkRemoveRouteCommand(cli),
		newNetworkAddSubnetCommand(cli),
		newNetworkRemoveSubnetCommand(cli),
		newNetworkAddLabelCommand(cli),
		newNetworkRemoveLabelCommand(cli),
		newNetworkEnableProtectionCommand(cli),
		newNetworkDisableProtectionCommand(cli),
	)
	return cmd
}
