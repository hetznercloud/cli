package cli

import "github.com/spf13/cobra"

func newContextCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "context [FLAGS]",
		Short:                 "Manage contexts",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runContext),
	}
	cmd.AddCommand(
		newContextCreateCommand(cli),
		newContextActiveCommand(cli),
		newContextActivateCommand(cli),
		newContextDeleteCommand(cli),
		newContextListCommand(cli),
	)
	return cmd
}

func runContext(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
