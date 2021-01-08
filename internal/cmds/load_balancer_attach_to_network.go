package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerAttachToNetworkCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "attach-to-network [FLAGS] LOADBALANCER",
		Short:                 "Attach a Load Balancer to a Network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerAttachToNetwork),
	}

	cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
	cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(cli.NetworkNames))
	cmd.MarkFlagRequired("network")

	cmd.Flags().IP("ip", nil, "IP address to assign to the Load Balancer (auto-assigned if omitted)")

	return cmd
}

func runLoadBalancerAttachToNetwork(cli *state.State, cmd *cobra.Command, args []string) error {
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
