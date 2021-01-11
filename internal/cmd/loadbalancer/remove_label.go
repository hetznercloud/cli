package loadbalancer

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newRemoveLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label [FLAGS] LOADBALANCER LABELKEY",
		Short: "Remove a label from a Load Balancer",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.LoadBalancerNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.LoadBalancerLabelKeys(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateLoadBalancerRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runLoadBalancerRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateLoadBalancerRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runLoadBalancerRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	labels := loadBalancer.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := loadBalancer.Labels[label]; !ok {
			return fmt.Errorf("label %s on Load Balancer %d does not exist", label, loadBalancer.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.LoadBalancerUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().LoadBalancer.Update(cli.Context, loadBalancer, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from Load Balancer %d\n", loadBalancer.ID)
	} else {
		fmt.Printf("Label %s removed from Load Balancer %d\n", args[1], loadBalancer.ID)
	}

	return nil
}
