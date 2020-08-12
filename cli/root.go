package cli

import (
	"time"

	"github.com/spf13/cobra"
)

func NewRootCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "hcloud",
		Short:                 "Hetzner Cloud CLI",
		Long:                  "A command-line interface for Hetzner Cloud",
		TraverseChildren:      true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		DisableFlagsInUseLine: true,
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
		newNetworkCommand(cli),
		newLoadBalancerCommand(cli),
		newLoadBalancerTypeCommand(cli),
		newCertificatesCommand(cli),
	)
	cmd.PersistentFlags().Duration("poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	return cmd
}
