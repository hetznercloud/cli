package volume

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.Volume, hcloud.VolumeChangeProtectionOpts]{
	ResourceNameSingular: "Volume",

	NameSuggestions: func(client hcapi2.Client) hcapi2.CompletionFunc {
		return client.Volume().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.VolumeChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.VolumeChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Volume, *hcloud.Response, error) {
		return s.Client().Volume().Get(cmd.Context(), idOrName)
	},

	ChangeProtectionFunction: func(ctx context.Context, s state.State, volume *hcloud.Volume, opts hcloud.VolumeChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().Volume().ChangeProtection(ctx, volume, opts)
	},

	IDOrName: func(volume *hcloud.Volume) string {
		return fmt.Sprint(volume.ID)
	},
}
