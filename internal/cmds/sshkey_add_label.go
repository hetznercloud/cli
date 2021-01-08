package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newSSHKeyAddLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] SSHKEY LABEL",
		Short:                 "Add a label to a SSH key",
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.SSHKeyNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateSSHKeyAddLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runSSHKeyAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateSSHKeyAddLabel(cmd *cobra.Command, args []string) error {
	label := splitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runSSHKeyAddLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]
	sshKey, _, err := cli.Client().SSHKey.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %s", idOrName)
	}
	label := splitLabel(args[1])

	if _, ok := sshKey.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on SSH key %d already exists", label[0], sshKey.ID)
	}
	labels := sshKey.Labels
	labels[label[0]] = label[1]
	opts := hcloud.SSHKeyUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().SSHKey.Update(cli.Context, sshKey, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to SSH key %d\n", label[0], sshKey.ID)

	return nil
}
