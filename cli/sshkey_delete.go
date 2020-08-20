package cli

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/spf13/cobra"
)

func newSSHKeyDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] SSHKEY",
		Short:                 "Delete a SSH key",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.SSHKeyNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runSSHKeyDelete),
	}
	return cmd
}

func runSSHKeyDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	sshKey, _, err := cli.Client().SSHKey.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %s", args[0])
	}

	_, err = cli.Client().SSHKey.Delete(cli.Context, sshKey)
	if err != nil {
		return err
	}

	fmt.Printf("SSH key %d deleted\n", sshKey.ID)
	return nil
}
