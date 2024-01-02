package loadbalancer

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DetachFromNetworkCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "detach-from-network [FLAGS] LOADBALANCER",
			Short:                 "Detach a Load Balancer from a Network",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.LoadBalancer().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
		cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(client.Network().Names))
		cmd.MarkFlagRequired("network")
		return cmd
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
		networkIDOrName, _ := cmd.Flags().GetString("network")
		network, _, err := client.Network().Get(ctx, networkIDOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", networkIDOrName)
		}

		opts := hcloud.LoadBalancerDetachFromNetworkOpts{
			Network: network,
		}
		action, _, err := client.LoadBalancer().DetachFromNetwork(ctx, loadBalancer, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(cmd, ctx, action); err != nil {
			return err
		}

		cmd.Printf("Load Balancer %d detached from Network %d\n", loadBalancer.ID, network.ID)
		return nil
	},
}
