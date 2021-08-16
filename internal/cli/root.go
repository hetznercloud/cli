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
	"github.com/hetznercloud/cli/internal/cmd/placementgroup"
	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/cmd/version"
	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewRootCommand(state *state.State, client hcapi2.Client) *cobra.Command {
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
		floatingip.NewCommand(state, client),
		image.NewCommand(state, client),
		server.NewCommand(state, client),
		sshkey.NewCommand(state, client),
		version.NewCommand(state),
		completion.NewCommand(state),
		servertype.NewCommand(state, client),
		context.NewCommand(state),
		datacenter.NewCommand(state, client),
		location.NewCommand(state, client),
		iso.NewCommand(state, client),
		volume.NewCommand(state, client),
		network.NewCommand(state, client),
		loadbalancer.NewCommand(state, client),
		loadbalancertype.NewCommand(state, client),
		certificate.NewCommand(state, client),
		firewall.NewCommand(state, client),
		placementgroup.NewCommand(state, client),
	)
	cmd.PersistentFlags().Duration("poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	return cmd
}
