package cli

import (
	"time"

	"github.com/spf13/cobra"
)

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
		newVolumeCommand(cli),
	)
	cmd.PersistentFlags().Duration("poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	return cmd
}

func runRoot(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
