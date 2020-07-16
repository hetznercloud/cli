package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLoadBalancerDeleteServiceCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete-service [FLAGS] LOADBALANCER",
		Short:                 "Deletes a service from a Load Balancer",
		Args:                  cobra.RangeArgs(1, 2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(cli.ensureToken),
		RunE:                  cli.wrap(runLoadBalancerDeleteService),
	}

	cmd.Flags().Int("listen-port", 0, "The listen port of the service you want to delete (required)")
	cmd.MarkFlagRequired("listen-port")
	return cmd
}

func runLoadBalancerDeleteService(cli *CLI, cmd *cobra.Command, args []string) error {
	listenPort, _ := cmd.Flags().GetInt("listen-port")
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}
	_, _, err = cli.Client().LoadBalancer.DeleteService(cli.Context, loadBalancer, listenPort)
	if err != nil {
		return err
	}

	fmt.Printf("Service on port %d deleted from Load Balancer %d\n", listenPort, loadBalancer.ID)
	return nil
}
