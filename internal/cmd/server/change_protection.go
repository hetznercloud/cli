package server

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.Server, hcloud.ServerChangeProtectionOpts]{
	ResourceNameSingular: "Server",

	NameSuggestions: func(client hcapi2.Client) hcapi2.CompletionFunc {
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

	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Server, *hcloud.Response, error) {
		return s.Client().Server().Get(cmd.Context(), idOrName)
	},

	ChangeProtectionFunction: func(ctx context.Context, s state.State, server *hcloud.Server, opts hcloud.ServerChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().Server().ChangeProtection(ctx, server, opts)
	},

	IDOrName: func(server *hcloud.Server) string {
		return fmt.Sprint(server.ID)
	},
}
