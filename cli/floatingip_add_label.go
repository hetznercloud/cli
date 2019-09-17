package cli

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPAddLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] FLOATINGIP LABEL",
		Short:                 "Add a label to a Floating IP",
		Args:                  cobra.ExactArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateFloatingIPAddLabel, cli.ensureToken),
		RunE:                  cli.wrap(runFloatingIPAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateFloatingIPAddLabel(cmd *cobra.Command, args []string) error {
	label := splitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runFloatingIPAddLabel(cli *CLI, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	idOrName := args[0]
	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %v", idOrName)
	}

	label := splitLabel(args[1])

	if _, ok := floatingIP.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on Floating IP %d already exists", label[0], floatingIP.ID)
	}
	labels := floatingIP.Labels
	labels[label[0]] = label[1]
	opts := hcloud.FloatingIPUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().FloatingIP.Update(cli.Context, floatingIP, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to Floating IP %d\n", label[0], floatingIP.ID)

	return nil
}
