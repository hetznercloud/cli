package cli

import "github.com/spf13/cobra"

func NewRootCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "hcloud",
		Short:            "Hetzner Cloud CLI",
		Long:             "A command-line interface for Hetzner Cloud",
		RunE:             cli.wrap(runRoot),
		TraverseChildren: true,
		SilenceUsage:     true,
		SilenceErrors:    true,
	}
	cmd.Flags().StringVar(&cli.Token, "token", "", "API token used for authentication")
	cmd.Flags().StringVar(&cli.Endpoint, "endpoint", Endpoint, "API endpoint URL")
	cmd.Flags().BoolVar(&cli.JSON, "json", false, "Output JSON API response")
	cmd.AddCommand(
		newFloatingIPCommand(cli),
		newServerCommand(cli),
		newSSHKeyCommand(cli),
		newVersionCommand(cli),
	)
	return cmd
}

func runRoot(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
