package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newSSHKeyDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "describe [flags] <id>",
		Short:            "Describe a SSH key",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runSSHKeyDescribe),
	}
	return cmd
}

func runSSHKeyDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid SSH key ID")
	}

	ctx := context.Background()

	sshKey, _, err := cli.Client().SSHKey.Get(ctx, id)
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %d", id)
	}

	fmt.Printf("ID:\t\t%d\n", sshKey.ID)
	fmt.Printf("Name:\t\t%s\n", sshKey.Name)
	fmt.Printf("Fingerprint:\t%s\n", sshKey.Fingerprint)
	fmt.Printf("Public Key:\n%s\n", strings.TrimSpace(sshKey.PublicKey))

	return nil
}
