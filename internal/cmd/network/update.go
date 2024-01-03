package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "Network",
	ShortDescription:     "Update a Network.\n\nTo enable or disable exposing routes to the vSwitch connection you can use the subcommand \"hcloud network expose-routes-to-vswitch\".",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Network().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Network name")
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		floatingIP := resource.(*hcloud.Network)
		updOpts := hcloud.NetworkUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Network().Update(s, floatingIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
