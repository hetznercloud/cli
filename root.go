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
		BashCompletionFunction: bashCompletionFunc,
	}
	cmd.Flags().BoolVar(&cli.JSON, "json", false, "Output JSON API response")
	cmd.AddCommand(
		newConfigureCommand(cli),
		newFloatingIPCommand(cli),
		newServerCommand(cli),
		newSSHKeyCommand(cli),
		newVersionCommand(cli),
		newCompletionCommand(cli),
	)
	return cmd
}

func runRoot(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
