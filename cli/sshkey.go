package cli

import "github.com/spf13/cobra"

func newSSHKeyCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh-key",
		Short:                 "Manage SSH keys",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runSSHKey),
	}
	cmd.AddCommand(
		newSSHKeyListCommand(cli),
		newSSHKeyCreateCommand(cli),
		newSSHKeyDeleteCommand(cli),
		newSSHKeyDescribeCommand(cli),
	)
	return cmd
}

func runSSHKey(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
