package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newLoadBalancerDisablePublicInterface(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disable-public-interface [FLAGS] LOADBALANCER",
		Short:                 "Disable the public interface of a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerDisablePublicInterface),
	}

	return cmd
}

func runLoadBalancerDisablePublicInterface(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	action, _, err := cli.Client().LoadBalancer.DisablePublicInterface(cli.Context, loadBalancer)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Public interface of Load Balancer %d was disabled\n", loadBalancer.ID)
	return nil
}
