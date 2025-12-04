package network

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.Network, hcloud.NetworkChangeProtectionOpts]{
	ResourceNameSingular: "Network",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.Network().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.NetworkChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.NetworkChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Network, *hcloud.Response, error) {
		return s.Client().Network().Get(s, idOrName)
	},

	ChangeProtectionFunction: func(s state.State, network *hcloud.Network, opts hcloud.NetworkChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().Network().ChangeProtection(s, network, opts)
	},

	GetIDFunction: func(network *hcloud.Network) int64 {
		return network.ID
	},
}
