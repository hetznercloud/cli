package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var ChangeTypeCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "change-type [FLAGS] LOADBALANCER LOADBALANCERTYPE",
			Short: "Change type of a Load Balancer",
			Args:  cobra.ExactArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.LoadBalancer().Names),
				cmpl.SuggestCandidatesF(client.LoadBalancerType().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		loadBalancer, _, err := client.LoadBalancer().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		loadBalancerTypeIDOrName := args[1]
		loadBalancerType, _, err := client.LoadBalancerType().Get(ctx, loadBalancerTypeIDOrName)
		if err != nil {
			return err
		}
		if loadBalancerType == nil {
			return fmt.Errorf("Load Balancer type not found: %s", loadBalancerTypeIDOrName)
		}

		opts := hcloud.LoadBalancerChangeTypeOpts{
			LoadBalancerType: loadBalancerType,
		}
		action, _, err := client.LoadBalancer().ChangeType(ctx, loadBalancer, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("LoadBalancer %d changed to type %s\n", loadBalancer.ID, loadBalancerType.Name)
		return nil
	},
}
