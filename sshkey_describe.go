package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
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

	sshKey, resp, err := cli.Client().SSHKey.Get(ctx, id)
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %d", id)
	}

	if cli.JSON {
		_, err = io.Copy(os.Stdout, resp.Body)
		return err
	}

	fmt.Printf("ID:\t\t%d\n", sshKey.ID)
	fmt.Printf("Name:\t\t%s\n", sshKey.Name)
	fmt.Printf("Fingerprint:\t%s\n", sshKey.Fingerprint)
	fmt.Printf("Public Key:\n%s\n", strings.TrimSpace(sshKey.PublicKey))

	return nil
}
