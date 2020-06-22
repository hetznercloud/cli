package cli

import "github.com/spf13/cobra"

func newLoadBalancerTypeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load-balancer-type",
		Short:                 "Manage Load Balancer types",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runLoadBalancerType),
	}
	cmd.AddCommand(
		newLoadBalancerTypenDescribeCommand(cli),
		newLoadBalancerTypeListCommand(cli),
	)
	return cmd
}

func runLoadBalancerType(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
