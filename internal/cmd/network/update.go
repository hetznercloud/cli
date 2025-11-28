package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.Network]{
	ResourceNameSingular: "Network",
	ShortDescription:     "Update a Network.\n\nTo enable or disable exposing routes to the vSwitch connection you can use the subcommand \"hcloud network expose-routes-to-vswitch\".",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Network, *hcloud.Response, error) {
		return s.Client().Network().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Network name")
	},
	Update: func(s state.State, _ *cobra.Command, network *hcloud.Network, flags map[string]pflag.Value) error {
		updOpts := hcloud.NetworkUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().Network().Update(s, network, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
