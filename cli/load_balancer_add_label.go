package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerAddLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] LOADBALANCER LABEL",
		Short:                 "Add a label to a Load Balancer",
		Args:                  cobra.ExactArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateLoadBalancerAddLabel, cli.ensureToken),
		RunE:                  cli.wrap(runLoadBalancerAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateLoadBalancerAddLabel(cmd *cobra.Command, args []string) error {
	label := splitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runLoadBalancerAddLabel(cli *CLI, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}
	label := splitLabel(args[1])

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
