package loadbalancer

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newAddLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] LOADBALANCER LABEL",
		Short:                 "Add a label to a Load Balancer",
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateLoadBalancerAddLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runLoadBalancerAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateLoadBalancerAddLabel(cmd *cobra.Command, args []string) error {
	label := util.SplitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runLoadBalancerAddLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}
	label := util.SplitLabel(args[1])

	if _, ok := loadBalancer.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on Load Balancer %d already exists", label[0], loadBalancer.ID)
	}
	labels := loadBalancer.Labels
	labels[label[0]] = label[1]
	opts := hcloud.LoadBalancerUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().LoadBalancer.Update(cli.Context, loadBalancer, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to Load Balancer %d\n", label[0], loadBalancer.ID)

	return nil
}
