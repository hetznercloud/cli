package cli

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerAttachToNetworkCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "attach-to-network [FLAGS] LOADBALANCER",
		Short:                 "Attach a Load Balancer to a Network",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerAttachToNetwork),
	}

	cmd.Flags().StringP("network", "n", "", "Network (ID or name)")
	cmd.Flag("network").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_network_names"},
	}
	cmd.MarkFlagRequired("network")

	cmd.Flags().IP("ip", nil, "IP address to assign to the Load Balancer (auto-assigned if omitted)")

	return cmd
}

func runLoadBalancerAttachToNetwork(cli *CLI, cmd *cobra.Command, args []string) error {
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

	ip, _ := cmd.Flags().GetIP("ip")

	opts := hcloud.LoadBalancerAttachToNetworkOpts{
		Network: network,
		IP:      ip,
	}
	action, _, err := cli.Client().LoadBalancer.AttachToNetwork(cli.Context, loadBalancer, opts)

	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Load Balancer %d attached to network %d\n", loadBalancer.ID, network.ID)
	return nil
}
