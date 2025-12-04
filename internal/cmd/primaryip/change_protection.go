package primaryip

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.PrimaryIP, hcloud.PrimaryIPChangeProtectionOpts]{
	ResourceNameSingular: "Primary IP",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.PrimaryIP().Names(false, false, nil)
	},

	ProtectionLevelsOptional: true,
	ProtectionLevels: map[string]func(opts *hcloud.PrimaryIPChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.PrimaryIPChangeProtectionOpts, value bool) {
			opts.Delete = value
		},
	},

	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.PrimaryIP, *hcloud.Response, error) {
		return s.Client().PrimaryIP().Get(s, idOrName)
	},

	ChangeProtectionFunction: func(s state.State, primaryIP *hcloud.PrimaryIP, opts hcloud.PrimaryIPChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		opts.ID = primaryIP.ID
		return s.Client().PrimaryIP().ChangeProtection(s, opts)
	},

	IDOrName: func(primaryIP *hcloud.PrimaryIP) string {
		return fmt.Sprint(primaryIP.ID)
	},
}
