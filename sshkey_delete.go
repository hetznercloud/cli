package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newSSHKeyDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] SSHKEY",
		Short:                 "Delete a SSH key",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runSSHKeyDelete),
	}
	return cmd
}

func runSSHKeyDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid SSH key ID")
	}
	sshKey := &hcloud.SSHKey{ID: id}

	_, err = cli.Client().SSHKey.Delete(cli.Context, sshKey)
	if err != nil {
		return err
	}

	fmt.Printf("SSH key %d deleted\n", id)
	return nil
}
