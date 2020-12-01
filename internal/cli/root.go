package cli

import (
	"time"

	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/cmd/completion"
	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/cmd/datacenter"
	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/cmd/iso"
	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/cmd/loadbalancertype"
	"github.com/hetznercloud/cli/internal/cmd/location"
	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/cmd/version"
	"github.com/hetznercloud/cli/internal/cmd/volume"
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
		floatingip.NewCommand(state),
		image.NewCommand(state),
		server.NewCommand(state),
		sshkey.NewCommand(state),
		version.NewCommand(state),
		completion.NewCommand(state),
		servertype.NewCommand(state),
		context.NewCommand(state),
		datacenter.NewCommand(state),
		location.NewCommand(state),
		iso.NewCommand(state),
		volume.NewCommand(state),
		network.NewCommand(state),
		loadbalancer.NewCommand(state),
		loadbalancertype.NewCommand(state),
		certificate.NewCommand(state),
		firewall.NewCommand(state),
	)
	cmd.PersistentFlags().Duration("poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	return cmd
}
