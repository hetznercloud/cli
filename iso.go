package cli

import "github.com/spf13/cobra"

func newISOCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "iso",
		Short:                 "Show information about ISOs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServer),
	}
	cmd.AddCommand(
		newISOListCommand(cli),
		newISODescribeCommand(cli),
	)
	return cmd
}

func runISO(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
