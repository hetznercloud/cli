package loadbalancertype

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load-balancer-type",
		Short:                 "Manage Load Balancer types",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newDescribeCommand(cli),
		newListCommand(cli),
	)
	return cmd
}
