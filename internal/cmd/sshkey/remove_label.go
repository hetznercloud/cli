package sshkey

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
		Use:   "remove-label [FLAGS] SSHKEY LABELKEY",
		Short: "Remove a label from a SSH key",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.SSHKeyNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.SSHKeyLabelKeys(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateSSHKeyRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runSSHKeyRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateSSHKeyRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runSSHKeyRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]
	sshKey, _, err := cli.Client().SSHKey.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %s", idOrName)
	}

	labels := sshKey.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := sshKey.Labels[label]; !ok {
			return fmt.Errorf("label %s on SSH key %d does not exist", label, sshKey.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.SSHKeyUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().SSHKey.Update(cli.Context, sshKey, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from SSH key %d\n", sshKey.ID)
	} else {
		fmt.Printf("Label %s removed from SSH key %d\n", args[1], sshKey.ID)
	}

	return nil
}
