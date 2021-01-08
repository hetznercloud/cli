package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewSSHKeyCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh-key",
		Short:                 "Manage SSH keys",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newSSHKeyListCommand(cli),
		newSSHKeyCreateCommand(cli),
		newSSHKeyUpdateCommand(cli),
		newSSHKeyDeleteCommand(cli),
		newSSHKeyDescribeCommand(cli),
		newSSHKeyAddLabelCommand(cli),
		newSSHKeyRemoveLabelCommand(cli),
	)
	return cmd
}
