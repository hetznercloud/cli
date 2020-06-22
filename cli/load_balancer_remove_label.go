package cli

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerRemoveLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-label [FLAGS] LOADBALANCER LABELKEY",
		Short:                 "Remove a label from a Load Balancer",
		Args:                  cobra.RangeArgs(1, 2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateLoadBalancerRemoveLabel, cli.ensureToken),
		RunE:                  cli.wrap(runLoadBalancerRemoveLabel),
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

func runLoadBalancerRemoveLabel(cli *CLI, cmd *cobra.Command, args []string) error {
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
