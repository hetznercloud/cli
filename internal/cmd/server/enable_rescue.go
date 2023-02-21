package server

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var EnableRescueCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "enable-rescue [FLAGS] SERVER",
			Short:                 "Enable rescue for a server",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "linux64", "Rescue type")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("linux64", "linux32"))

		cmd.Flags().StringSlice("ssh-key", nil, "ID or name of SSH key to inject (can be specified multiple times)")
		cmd.RegisterFlagCompletionFunc("ssh-key", cmpl.SuggestCandidatesF(client.SSHKey().Names))
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		var (
			opts hcloud.ServerEnableRescueOpts
		)
		rescueType, _ := cmd.Flags().GetString("type")
		opts.Type = hcloud.ServerRescueType(rescueType)

		sshKeys, _ := cmd.Flags().GetStringSlice("ssh-key")
		for _, sshKeyIDOrName := range sshKeys {
			sshKey, _, err := client.SSHKey().Get(ctx, sshKeyIDOrName)
			if err != nil {
				return err
			}
			if sshKey == nil {
				return fmt.Errorf("SSH key not found: %s", sshKeyIDOrName)
			}
			opts.SSHKeys = append(opts.SSHKeys, sshKey)
		}

		result, _, err := client.Server().EnableRescue(ctx, server, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, result.Action); err != nil {
			return err
		}

		fmt.Printf("Rescue enabled for server %d with root password: %s\n", server.ID, result.RootPassword)
		return nil
	},
}
