package loadbalancer

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newDeleteCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] LOADBALANCER",
		Short:                 "Delete a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDelete),
	}
	return cmd
}

func runDelete(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load balancer not found: %s", idOrName)
	}

	_, err = cli.Client().LoadBalancer.Delete(cli.Context, loadBalancer)
	if err != nil {
		return err
	}

	fmt.Printf("Load Balancer %d deleted\n", loadBalancer.ID)
	return nil
}
