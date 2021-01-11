package network

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
		Use:                   "add-label [FLAGS] NETWORK LABEL",
		Short:                 "Add a label to a network",
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.NetworkNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateNetworkAddLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runNetworkAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateNetworkAddLabel(cmd *cobra.Command, args []string) error {
	label := util.SplitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runNetworkAddLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]
	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}
	label := util.SplitLabel(args[1])

	if _, ok := network.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on network %d already exists", label[0], network.ID)
	}
	labels := network.Labels
	labels[label[0]] = label[1]
	opts := hcloud.NetworkUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Network.Update(cli.Context, network, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to network %d\n", label[0], network.ID)

	return nil
}
