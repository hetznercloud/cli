package cli

import "github.com/spf13/cobra"

func newFloatingIPCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "floatingip",
		Short:            "Manage Floating IPs",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE:             cli.wrap(runFloatingIP),
	}
	cmd.AddCommand(
		newFloatingIPListCommand(cli),
	)
	return cmd
}

func runFloatingIP(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
