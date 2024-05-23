package e2e_tests

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cli"
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
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/cmd/version"
	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func TestMain(m *testing.M) {
	tok := os.Getenv("HCLOUD_TOKEN")
	if tok == "" {
		fmt.Println("HCLOUD_TOKEN is not set")
		os.Exit(1)
		return
	}
	os.Exit(m.Run())
}

func newRootCommand(t *testing.T) *cobra.Command {
	cfg := config.NewConfig()
	if err := config.ReadConfig(cfg, "config.toml"); err != nil {
		t.Fatalf("unable to read config file \"%s\": %s\n", cfg.Path(), err)
	}

	s, err := state.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	rootCommand := cli.NewRootCommand(s)

	util.AddGroup(rootCommand, "resource", "Resources",
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
	)

	rootCommand.AddCommand(
		version.NewCommand(s),
		completion.NewCommand(s),
		context.NewCommand(s),
		configCmd.NewCommand(s),
	)

	return rootCommand
}

func runCommand(t *testing.T, args ...string) (string, error) {
	cmd := newRootCommand(t)
	var buf bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetOut(&buf)
	err := cmd.Execute()
	return buf.String(), err
}
