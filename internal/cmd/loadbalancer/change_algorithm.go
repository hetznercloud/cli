package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var ChangeAlgorithmCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "change-algorithm LOADBALANCER FLAGS",
			Short:                 "Changes the algorithm of a Load Balancer",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.LoadBalancer().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("algorithm-type", "", "The new algorithm of the Load Balancer")
		cmd.RegisterFlagCompletionFunc("algorithm-type", cmpl.SuggestCandidates(
			string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
			string(hcloud.LoadBalancerAlgorithmTypeLeastConnections),
		))
		cmd.MarkFlagRequired("algorithm-type")

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		algorithm, _ := cmd.Flags().GetString("algorithm-type")
		loadBalancer, _, err := client.LoadBalancer().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		action, _, err := client.LoadBalancer().ChangeAlgorithm(ctx, loadBalancer, hcloud.LoadBalancerChangeAlgorithmOpts{Type: hcloud.LoadBalancerAlgorithmType(algorithm)})
		if err != nil {
			return err
		}
		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}
		fmt.Printf("Algorithm for Load Balancer %d was changed\n", loadBalancer.ID)

		return nil
	},
}
