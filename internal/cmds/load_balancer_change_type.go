package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerChangeTypeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-type [FLAGS] LOADBALANCER LOADBALANCERTYPE",
		Short: "Change type of a Load Balancer",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.LoadBalancerNames),
			cmpl.SuggestCandidatesF(cli.LoadBalancerTypeNames),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerChangeType),
	}

	return cmd
}

func runLoadBalancerChangeType(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	loadBalancerTypeIDOrName := args[1]
	loadBalancerType, _, err := cli.Client().LoadBalancerType.Get(cli.Context, loadBalancerTypeIDOrName)
	if err != nil {
		return err
	}
	if loadBalancerType == nil {
		return fmt.Errorf("Load Balancer type not found: %s", loadBalancerTypeIDOrName)
	}

	opts := hcloud.LoadBalancerChangeTypeOpts{
		LoadBalancerType: loadBalancerType,
	}
	action, _, err := cli.Client().LoadBalancer.ChangeType(cli.Context, loadBalancer, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("LoadBalancer %d changed to type %s\n", loadBalancer.ID, loadBalancerType.Name)
	return nil
}
