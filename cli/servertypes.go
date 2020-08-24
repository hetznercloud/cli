package cli

import "github.com/spf13/cobra"

func newServerTypeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "server-type",
		Short:                 "Manage server types",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newServerTypeListCommand(cli),
		newServerTypeDescribeCommand(cli),
	)
	return cmd
}
