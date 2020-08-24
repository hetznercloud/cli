package cli

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerChangeAlgorithmCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "change-algorithm LOADBALANCER FLAGS",
		Short:                 "Changes the algorithm of a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerChangeAlgorithm),
	}

	cmd.Flags().String("algorithm-type", "", "The new algorithm of the Load Balancer")
	cmd.RegisterFlagCompletionFunc("algorithm-type", cmpl.SuggestCandidates(
		string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
		string(hcloud.LoadBalancerAlgorithmTypeRoundRobin),
	))
	cmd.MarkFlagRequired("algorithm-type")

	return cmd
}

func runLoadBalancerChangeAlgorithm(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	algorithm, _ := cmd.Flags().GetString("algorithm-type")
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	action, _, err := cli.Client().LoadBalancer.ChangeAlgorithm(cli.Context, loadBalancer, hcloud.LoadBalancerChangeAlgorithmOpts{Type: hcloud.LoadBalancerAlgorithmType(algorithm)})
	if err != nil {
		return err
	}
	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Algorithm for Load Balancer %d was changed\n", loadBalancer.ID)

	return nil
}
