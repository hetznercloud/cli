package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"

	"github.com/spf13/cobra"
)

func newLoadBalancerDetachFromNetworkCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "detach-from-network [FLAGS] LOADBALANCER",
		Short:                 "Detach a Load Balancer from a Network",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerDetachFromNetwork),
	}
	cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
	cmd.Flag("network").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_network_names"},
	}
	cmd.MarkFlagRequired("network")
	return cmd
}

func runLoadBalancerDetachFromNetwork(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}
	networkIDOrName, _ := cmd.Flags().GetString("network")
	network, _, err := cli.Client().Network.Get(cli.Context, networkIDOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", networkIDOrName)
	}

	opts := hcloud.LoadBalancerDetachFromNetworkOpts{
		Network: network,
	}
	action, _, err := cli.Client().LoadBalancer.DetachFromNetwork(cli.Context, loadBalancer, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Load Balancer %d detached from Network %d\n", loadBalancer.ID, network.ID)
	return nil
}
