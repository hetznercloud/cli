package firewall

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
		Use:   "remove-label [FLAGS] FIREWALL LABELKEY",
		Short: "Remove a label from a Firewall",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.FirewallNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.FirewallLabelKeys(args[0])
			})),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateFirewallRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runFirewallRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateFirewallRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runFirewallRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	idOrName := args[0]
	firewall, _, err := cli.Client().Firewall.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if firewall == nil {
		return fmt.Errorf("Firewall not found: %v", idOrName)
	}

	labels := firewall.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := firewall.Labels[label]; !ok {
			return fmt.Errorf("label %s on Firewall %d does not exist", label, firewall.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.FirewallUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Firewall.Update(cli.Context, firewall, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from Firewall %d\n", firewall.ID)
	} else {
		fmt.Printf("Label %s removed from Firewall %d\n", args[1], firewall.ID)
	}

	return nil
}
