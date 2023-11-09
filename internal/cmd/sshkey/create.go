package sshkey

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create FLAGS",
			Short: "Create a SSH key",
			Args:  cobra.NoArgs,
		}
		cmd.Flags().String("name", "", "Key name (required)")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().String("public-key", "", "Public key")
		cmd.Flags().String("public-key-from-file", "", "Path to file containing public key")
		cmd.MarkFlagsMutuallyExclusive("public-key", "public-key-from-file")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
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
				data, err = io.ReadAll(os.Stdin)
			} else {
				data, err = os.ReadFile(publicKeyFile)
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
		sshKey, _, err := client.SSHKey().Create(ctx, opts)
		if err != nil {
			return err
		}

		cmd.Printf("SSH key %d created\n", sshKey.ID)

		return nil
	},
}
