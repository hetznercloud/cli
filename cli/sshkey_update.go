package cli

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newSSHKeyUpdateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] SSHKEY",
		Short:                 "Update a SSH Key",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runSSHKeyUpdate),
	}

	cmd.Flags().String("name", "", "SSH Key name")

	return cmd
}

func runSSHKeyUpdate(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	sshKey, _, err := cli.Client().SSHKey.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %s", idOrName)
	}

	name, _ := cmd.Flags().GetString("name")
	opts := hcloud.SSHKeyUpdateOpts{
		Name: name,
	}
	if opts.Name == "" {
		return errors.New("no updates")
	}

	_, _, err = cli.Client().SSHKey.Update(cli.Context, sshKey, opts)
	if err != nil {
		return err
	}
	fmt.Printf("SSH Key %s updated\n", sshKey.Name)
	return nil
}
