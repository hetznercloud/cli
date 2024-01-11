package main

import (
	"log"
	"os"

	"github.com/hetznercloud/cli/internal/cli"
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
	"github.com/hetznercloud/cli/internal/state/config"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("hcloud: ")
	log.SetOutput(os.Stderr)
}

func main() {
	configPath := os.Getenv("HCLOUD_CONFIG")
	if configPath == "" {
		configPath = config.DefaultConfigPath()
	}

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("unable to read config file %q: %s\n", configPath, err)
	}

	s, err := state.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	rootCommand := cli.NewRootCommand(s)
	rootCommand.AddCommand(
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

	if err := rootCommand.Execute(); err != nil {
		log.Fatalln(err)
	}
}
