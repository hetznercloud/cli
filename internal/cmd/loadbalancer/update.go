package loadbalancer

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newUpdateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] LOADBALANCER",
		Short:                 "Update a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerUpdate),
	}

	cmd.Flags().String("name", "", "Load Balancer name")

	return cmd
}

func runLoadBalancerUpdate(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	name, _ := cmd.Flags().GetString("name")
	opts := hcloud.LoadBalancerUpdateOpts{
		Name: name,
	}
	if opts.Name == "" {
		return errors.New("no updates")
	}

	_, _, err = cli.Client().LoadBalancer.Update(cli.Context, loadBalancer, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Load Balancer %d updated\n", loadBalancer.ID)
	return nil
}
