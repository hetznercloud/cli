package sshkey

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh-key",
		Short:                 "Manage SSH keys",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newListCommand(cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		newDeleteCommand(cli),
		newDescribeCommand(cli),
		newAddLabelCommand(cli),
		newRemoveLabelCommand(cli),
	)
	return cmd
}
