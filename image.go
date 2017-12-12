package cli

import "github.com/spf13/cobra"

func newImageCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "image",
		Short:            "Manage Images",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE:             cli.wrap(runImage),
	}
	cmd.AddCommand(
		newImageListCommand(cli),
		newImageDeleteCommand(cli),
		newImageDescribeCommand(cli),
	)
	return cmd
}

func runImage(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
