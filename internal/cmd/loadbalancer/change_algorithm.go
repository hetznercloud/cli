package loadbalancer

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeAlgorithmCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "change-algorithm --algorithm-type <round_robin|least_connections> <load-balancer>",
			Short:                 "Changes the algorithm of a Load Balancer",
			Args:                  util.Validate,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.LoadBalancer().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("algorithm-type", "", "New Load Balancer algorithm (round_robin, least_connections) (required)")
		cmd.RegisterFlagCompletionFunc("algorithm-type", cmpl.SuggestCandidates(
			string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
			string(hcloud.LoadBalancerAlgorithmTypeLeastConnections),
		))
		cmd.MarkFlagRequired("algorithm-type")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		algorithm, _ := cmd.Flags().GetString("algorithm-type")
		loadBalancer, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		action, _, err := s.Client().LoadBalancer().ChangeAlgorithm(s, loadBalancer, hcloud.LoadBalancerChangeAlgorithmOpts{Type: hcloud.LoadBalancerAlgorithmType(algorithm)})
		if err != nil {
			return err
		}
		if err := s.ActionProgress(cmd, s, action); err != nil {
			return err
		}
		cmd.Printf("Algorithm for Load Balancer %d was changed\n", loadBalancer.ID)

		return nil
	},
}
