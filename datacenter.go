package cli

import "github.com/spf13/cobra"

func newDatacenterCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "datacenter",
		Short:                 "Show information about datacenters",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runDatacenter),
	}
	cmd.AddCommand(
		newDatacenterListCommand(cli),
		newDatacenterDescribeCommand(cli),
	)
	return cmd
}

func runDatacenter(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
