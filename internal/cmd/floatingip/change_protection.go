package floatingip

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.FloatingIP, hcloud.FloatingIPChangeProtectionOpts]{
	ResourceNameSingular: "Floating IP",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.FloatingIP().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.FloatingIPChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.FloatingIPChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.FloatingIP, *hcloud.Response, error) {
		return s.Client().FloatingIP().Get(s, idOrName)
	},

	ChangeProtectionFunction: func(s state.State, floatingIP *hcloud.FloatingIP, opts hcloud.FloatingIPChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().FloatingIP().ChangeProtection(s, floatingIP, opts)
	},

	IDOrName: func(floatingIP *hcloud.FloatingIP) string {
		return fmt.Sprint(floatingIP.ID)
	},
}
