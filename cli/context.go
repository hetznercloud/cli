package cli

import "github.com/spf13/cobra"

func newContextCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "context [FLAGS]",
		Short:                 "Manage contexts",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newContextCreateCommand(cli),
		newContextActiveCommand(cli),
		newContextUseCommand(cli),
		newContextDeleteCommand(cli),
		newContextListCommand(cli),
	)
	return cmd
}
