package server

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
		Use:   "remove-label [FLAGS] SERVER LABELKEY",
		Short: "Remove a label from a server",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.ServerNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.ServerLabelKeys(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	labels := server.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := server.Labels[label]; !ok {
			return fmt.Errorf("label %s on server %d does not exist", label, server.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.ServerUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Server.Update(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from server %d\n", server.ID)
	} else {
		fmt.Printf("Label %s removed from server %d\n", args[1], server.ID)
	}

	return nil
}
