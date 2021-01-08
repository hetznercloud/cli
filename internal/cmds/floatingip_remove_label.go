package cmds

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPRemoveLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label [FLAGS] FLOATINGIP LABELKEY",
		Short: "Remove a label from a Floating IP",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.FloatingIPNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.FloatingIPLabelKeys(args[0])
			})),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateFloatingIPRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runFloatingIPRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateFloatingIPRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runFloatingIPRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	idOrName := args[0]
	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %v", idOrName)
	}

	labels := floatingIP.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := floatingIP.Labels[label]; !ok {
			return fmt.Errorf("label %s on Floating IP %d does not exist", label, floatingIP.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.FloatingIPUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().FloatingIP.Update(cli.Context, floatingIP, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from Floating IP %d\n", floatingIP.ID)
	} else {
		fmt.Printf("Label %s removed from Floating IP %d\n", args[1], floatingIP.ID)
	}

	return nil
}
