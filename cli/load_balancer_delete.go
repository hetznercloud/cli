package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLoadBalancerDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] LOADBALANCER",
		Short:                 "Delete a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerDelete),
	}
	return cmd
}

func runLoadBalancerDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load balancer not found: %s", idOrName)
	}

	_, err = cli.Client().LoadBalancer.Delete(cli.Context, loadBalancer)
	if err != nil {
		return err
	}

	fmt.Printf("Load Balancer %d deleted\n", loadBalancer.ID)
	return nil
}
