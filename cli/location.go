package cli

import "github.com/spf13/cobra"

func newLocationCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "location",
		Short:                 "Manage locations",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newLocationListCommand(cli),
		newLocationDescribeCommand(cli),
	)
	return cmd
}
