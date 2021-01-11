package server

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newEnableRescueCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "enable-rescue [FLAGS] SERVER",
		Short:                 "Enable rescue for a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerEnableRescue),
	}
	cmd.Flags().String("type", "linux64", "Rescue type")
	cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("linux64", "linux32", "freebsd64"))

	cmd.Flags().StringSlice("ssh-key", nil, "ID or name of SSH key to inject (can be specified multiple times)")
	cmd.RegisterFlagCompletionFunc("ssh-key", cmpl.SuggestCandidatesF(cli.SSHKeyNames))
	return cmd
}

func runServerEnableRescue(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
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
		sshKey, _, err := cli.Client().SSHKey.Get(cli.Context, sshKeyIDOrName)
		if err != nil {
			return err
		}
		if sshKey == nil {
			return fmt.Errorf("SSH key not found: %s", sshKeyIDOrName)
		}
		opts.SSHKeys = append(opts.SSHKeys, sshKey)
	}

	result, _, err := cli.Client().Server.EnableRescue(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}

	fmt.Printf("Rescue enabled for server %d with root password: %s\n", server.ID, result.RootPassword)
	return nil
}
