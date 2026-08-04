package storagebox

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.StorageBox, hcloud.StorageBoxChangeProtectionOpts]{
	ResourceNameSingular: "Storage Box",

	NameSuggestions: func(client hcapi2.Client) hcapi2.CompletionFunc {
		return client.StorageBox().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.StorageBoxChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.StorageBoxChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.StorageBox, *hcloud.Response, error) {
		return s.Client().StorageBox().Get(cmd.Context(), idOrName)
	},

	ChangeProtectionFunction: func(ctx context.Context, s state.State, storageBox *hcloud.StorageBox, opts hcloud.StorageBoxChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().StorageBox().ChangeProtection(ctx, storageBox, opts)
	},

	IDOrName: func(storageBox *hcloud.StorageBox) string {
		return fmt.Sprint(storageBox.ID)
	},
}
