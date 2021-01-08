package cli

import (
	"time"

	"github.com/hetznercloud/cli/internal/cmds"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewRootCommand(state *state.State) *cobra.Command {
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
		cmds.NewFloatingIPCommand(state),
		cmds.NewImageCommand(state),
		cmds.NewServerCommand(state),
		cmds.NewSSHKeyCommand(state),
		cmds.NewVersionCommand(state),
		cmds.NewCompletionCommand(state),
		cmds.NewServerTypeCommand(state),
		cmds.NewContextCommand(state),
		cmds.NewDatacenterCommand(state),
		cmds.NewLocationCommand(state),
		cmds.NewISOCommand(state),
		cmds.NewVolumeCommand(state),
		cmds.NewNetworkCommand(state),
		cmds.NewLoadBalancerCommand(state),
		cmds.NewLoadBalancerTypeCommand(state),
		cmds.NewCertificatesCommand(state),
	)
	cmd.PersistentFlags().Duration("poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	return cmd
}
