package cmds

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkRemoveLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label [FLAGS] NETWORK LABELKEY",
		Short: "Remove a label from a network",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.NetworkNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.NetworkLabelKeys(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateNetworkRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runNetworkRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateNetworkRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runNetworkRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]
	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}

	labels := network.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := network.Labels[label]; !ok {
			return fmt.Errorf("label %s on network %d does not exist", label, network.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.NetworkUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Network.Update(cli.Context, network, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from network %d\n", network.ID)
	} else {
		fmt.Printf("Label %s removed from network %d\n", args[1], network.ID)
	}

	return nil
}
