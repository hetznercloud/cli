package cli

import "github.com/spf13/cobra"

func newLocationCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "location",
		Short:                 "Manage locations",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runServer),
	}
	cmd.AddCommand(
		newLocationListCommand(cli),
		newLocationDescribeCommand(cli),
	)
	return cmd
}

func runLocation(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
