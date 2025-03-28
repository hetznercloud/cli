package loadbalancertype

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load-balancer-type",
		Short:                 "View Load Balancer Types",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		DescribeCmd.CobraCommand(s),
		ListCmd.CobraCommand(s),
	)
	return cmd
}
