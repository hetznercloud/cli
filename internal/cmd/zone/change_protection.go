package zone

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.Zone, hcloud.ZoneChangeProtectionOpts]{
	ResourceNameSingular: "Zone",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.Zone().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.ZoneChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.ZoneChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Zone, *hcloud.Response, error) {
		return s.Client().Zone().Get(s, idOrName)
	},

	ChangeProtectionFunction: func(s state.State, zone *hcloud.Zone, opts hcloud.ZoneChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().Zone().ChangeProtection(s, zone, opts)
	},

	IDOrName: func(zone *hcloud.Zone) string {
		return zone.Name
	},
}
