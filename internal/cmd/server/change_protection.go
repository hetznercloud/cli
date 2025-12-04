package server

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.Server, hcloud.ServerChangeProtectionOpts]{
	ResourceNameSingular: "Server",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.Server().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.ServerChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.ServerChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
		"rebuild": func(opts *hcloud.ServerChangeProtectionOpts, value bool) {
			opts.Rebuild = &value
		},
	},

	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Server, *hcloud.Response, error) {
		return s.Client().Server().Get(s, idOrName)
	},

	ChangeProtectionFunction: func(s state.State, server *hcloud.Server, opts hcloud.ServerChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().Server().ChangeProtection(s, server, opts)
	},

	GetIDFunction: func(server *hcloud.Server) int64 {
		return server.ID
	},
}
