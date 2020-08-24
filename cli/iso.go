package cli

import "github.com/spf13/cobra"

func newISOCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "iso",
		Short:                 "Manage ISOs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newISOListCommand(cli),
		newISODescribeCommand(cli),
	)
	return cmd
}
