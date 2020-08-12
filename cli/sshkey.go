package cli

import "github.com/spf13/cobra"

func newSSHKeyCommand(cli *CLI) *cobra.Command {
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
