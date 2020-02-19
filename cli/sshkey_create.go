package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newSSHKeyCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create FLAGS",
		Short:                 "Create a SSH key",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateSSHKeyCreate, cli.ensureToken),
		RunE:                  cli.wrap(runSSHKeyCreate),
	}
	cmd.Flags().String("name", "", "Key name")
	cmd.Flags().String("public-key", "", "Public key")
	cmd.Flags().String("public-key-from-file", "", "Path to file containing public key")
	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")
	return cmd
}

func validateSSHKeyCreate(cmd *cobra.Command, args []string) error {
	if name, _ := cmd.Flags().GetString("name"); name == "" {
		return errors.New("flag --name is required")
	}

	publicKey, _ := cmd.Flags().GetString("public-key")
	publicKeyFile, _ := cmd.Flags().GetString("public-key-from-file")
	if publicKey != "" && publicKeyFile != "" {
		return errors.New("flags --public-key and --public-key-from-file are mutually exclusive")
	}

	return nil
}

func runSSHKeyCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	publicKey, _ := cmd.Flags().GetString("public-key")
	publicKeyFile, _ := cmd.Flags().GetString("public-key-from-file")
	labels, _ := cmd.Flags().GetStringToString("label")

	if publicKeyFile != "" {
		var (
			data []byte
			err  error
		)
		if publicKeyFile == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(publicKeyFile)
		}
		if err != nil {
			return err
		}
		publicKey = string(data)
	}

	opts := hcloud.SSHKeyCreateOpts{
		Name:      name,
		PublicKey: publicKey,
		Labels:    labels,
	}
	sshKey, _, err := cli.Client().SSHKey.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	fmt.Printf("SSH key %d created\n", sshKey.ID)

	return nil
}
