package cli

import "github.com/spf13/cobra"

func newServerTypeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "servertype",
		Short:            "Show information about servertypes",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE:             cli.wrap(runServer),
	}
	cmd.AddCommand(
		newServerTypeListCommand(cli),
		newServerTypeDescribeCommand(cli),
	)
	return cmd
}

func runServerType(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
