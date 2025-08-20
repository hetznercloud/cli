package cli

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/all"
	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/cmd/completion"
	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
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
	"github.com/hetznercloud/cli/internal/cmd/storageboxtype"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/cmd/version"
	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
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

	util.AddGroup(cmd, "resource", "Resources",
		all.NewCommand(s),
		floatingip.NewCommand(s),
		image.NewCommand(s),
		server.NewCommand(s),
		sshkey.NewCommand(s),
		servertype.NewCommand(s),
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
		storageboxtype.NewCommand(s),
	)

	cmd.AddCommand(
		version.NewCommand(s),
		completion.NewCommand(s),
		context.NewCommand(s),
		configCmd.NewCommand(s),
	)

	cmd.PersistentFlags().AddFlagSet(s.Config().FlagSet())

	for _, opt := range config.Options {
		f := opt.GetFlagCompletionFunc()
		if !opt.HasFlags(config.OptionFlagPFlag) || f == nil {
			continue
		}
		// opt.FlagName() is prefixed with --
		flagName := opt.FlagName()[2:]
		_ = cmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return f(s.Client(), s.Config(), cmd, args, toComplete)
		})
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		out := cmd.OutOrStdout()
		quiet, err := config.OptionQuiet.Get(s.Config())
		if err != nil {
			return err
		}
		if quiet {
			// We save the original output in cmd.errWriter so that we can still use it if we need it later.
			cmd.SetErr(out)
			out = io.Discard
		}
		cmd.SetOut(out)
		return nil
	}
	return cmd
}
