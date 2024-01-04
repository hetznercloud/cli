package cli

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/all"
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
	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/cmd/version"
	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func NewRootCommand(s state.State) *cobra.Command {
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
		all.NewCommand(s),
		floatingip.NewCommand(s),
		image.NewCommand(s),
		server.NewCommand(s),
		sshkey.NewCommand(s),
		version.NewCommand(s),
		completion.NewCommand(s),
		servertype.NewCommand(s),
		context.NewCommand(s),
		datacenter.NewCommand(s),
		location.NewCommand(s),
		iso.NewCommand(s),
		volume.NewCommand(s),
		network.NewCommand(s),
		loadbalancer.NewCommand(s),
		loadbalancertype.NewCommand(s),
		certificate.NewCommand(s),
		firewall.NewCommand(s),
		placementgroup.NewCommand(s),
		primaryip.NewCommand(s),
	)
	cmd.PersistentFlags().Duration("poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	cmd.SetOut(os.Stdout)
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		pollInterval, err := cmd.Flags().GetDuration("poll-interval")
		if err != nil {
			return err
		}
		s.Client().WithOpts(hcloud.WithPollBackoffFunc(hcloud.ConstantBackoff(pollInterval)))
		return nil
	}
	return cmd
}
