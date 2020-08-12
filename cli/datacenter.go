package cli

import "github.com/spf13/cobra"

func newDatacenterCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "datacenter",
		Short:                 "Manage datacenters",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newDatacenterListCommand(cli),
		newDatacenterDescribeCommand(cli),
	)
	return cmd
}
