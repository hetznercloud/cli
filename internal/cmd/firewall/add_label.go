package firewall

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func newAddLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] FIREWALL LABEL",
		Short:                 "Add a label to a Firewall",
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FirewallNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateFirewallAddLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runFirewallAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateFirewallAddLabel(cmd *cobra.Command, args []string) error {
	label := util.SplitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runFirewallAddLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	idOrName := args[0]
	firewall, _, err := cli.Client().Firewall.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if firewall == nil {
		return fmt.Errorf("Firewall not found: %v", idOrName)
	}

	label := util.SplitLabel(args[1])

	if _, ok := firewall.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on Firewall %d already exists", label[0], firewall.ID)
	}
	labels := firewall.Labels
	labels[label[0]] = label[1]
	opts := hcloud.FirewallUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Firewall.Update(cli.Context, firewall, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to Firewall %d\n", label[0], firewall.ID)

	return nil
}
