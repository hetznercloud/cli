package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var EnableRescueCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "enable-rescue [options] <server>",
			Short:                 "Enable rescue for a server",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("type", "linux64", "Rescue type")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("linux64"))

		cmd.Flags().StringSlice("ssh-key", nil, "ID or name of SSH Key to inject (can be specified multiple times)")
		_ = cmd.RegisterFlagCompletionFunc("ssh-key", cmpl.SuggestCandidatesF(client.SSHKey().Names))
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}

		var (
			opts hcloud.ServerEnableRescueOpts
		)
		rescueType, _ := cmd.Flags().GetString("type")

		opts.Type = hcloud.ServerRescueType(rescueType)
		switch opts.Type {
		case hcloud.ServerRescueTypeLinux64:
			break
		case hcloud.ServerRescueTypeLinux32: //nolint:staticcheck
			return fmt.Errorf("rescue type not supported anymore: %s", opts.Type)
		default:
			return fmt.Errorf("invalid rescue type: %s", opts.Type)
		}

		sshKeys, _ := cmd.Flags().GetStringSlice("ssh-key")
		for _, sshKeyIDOrName := range sshKeys {
			sshKey, _, err := s.Client().SSHKey().Get(s, sshKeyIDOrName)
			if err != nil {
				return err
			}
			if sshKey == nil {
				return fmt.Errorf("SSH Key not found: %s", sshKeyIDOrName)
			}
			opts.SSHKeys = append(opts.SSHKeys, sshKey)
		}

		result, _, err := s.Client().Server().EnableRescue(s, server, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return err
		}

		cmd.Printf("Rescue enabled for Server %d with root password: %s\n", server.ID, result.RootPassword)
		return nil
	},
}
