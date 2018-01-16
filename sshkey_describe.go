package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newSSHKeyDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SSHKEY",
		Short:                 "Describe a SSH key",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureActiveContext,
		RunE:                  cli.wrap(runSSHKeyDescribe),
	}
	return cmd
}

func runSSHKeyDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	sshKey, _, err := cli.Client().SSHKey.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %s", args[0])
	}

	fmt.Printf("ID:\t\t%d\n", sshKey.ID)
	fmt.Printf("Name:\t\t%s\n", sshKey.Name)
	fmt.Printf("Fingerprint:\t%s\n", sshKey.Fingerprint)
	fmt.Printf("Public Key:\n%s\n", strings.TrimSpace(sshKey.PublicKey))

	return nil
}
