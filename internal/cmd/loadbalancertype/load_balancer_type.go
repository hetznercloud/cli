package loadbalancertype

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load-balancer-type",
		Short:                 "Manage Load Balancer types",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		ListCmd.CobraCommand(cli.Context, client, cli),
	)
	return cmd
}
