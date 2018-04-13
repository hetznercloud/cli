package cli

import "github.com/spf13/cobra"

func NewRootCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "hcloud",
		Short:                  "Hetzner Cloud CLI",
		Long:                   "A command-line interface for Hetzner Cloud",
		RunE:                   cli.wrap(runRoot),
		TraverseChildren:       true,
		SilenceUsage:           true,
		SilenceErrors:          true,
		DisableFlagsInUseLine:  true,
		BashCompletionFunction: bashCompletionFunc,
	}
	cmd.AddCommand(
		newFloatingIPCommand(cli),
		newImageCommand(cli),
		newServerCommand(cli),
		newSSHKeyCommand(cli),
		newVersionCommand(cli),
		newCompletionCommand(cli),
		newServerTypeCommand(cli),
		newContextCommand(cli),
		newDatacenterCommand(cli),
		newLocationCommand(cli),
		newISOCommand(cli),
	)

	cmd.PersistentFlags().StringArrayP("output", "o", []string{}, "Output options. One of: noheader|columns=[col1,col2]")

	return cmd
}

func runRoot(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
