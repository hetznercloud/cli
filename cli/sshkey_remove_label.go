package cli

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newSSHKeyRemoveLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-label [FLAGS] SSHKEY LABELKEY",
		Short:                 "Remove a label from a SSH key",
		Args:                  cobra.RangeArgs(1, 2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateSSHKeyRemoveLabel, cli.ensureToken),
		RunE:                  cli.wrap(runSSHKeyRemoveLabel),
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

func runSSHKeyRemoveLabel(cli *CLI, cmd *cobra.Command, args []string) error {
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
