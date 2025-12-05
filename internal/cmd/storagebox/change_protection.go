package storagebox

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.StorageBox, hcloud.StorageBoxChangeProtectionOpts]{
	ResourceNameSingular: "Storage Box",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.StorageBox().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.StorageBoxChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.StorageBoxChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.StorageBox, *hcloud.Response, error) {
		return s.Client().StorageBox().Get(s, idOrName)
	},

	ChangeProtectionFunction: func(s state.State, storageBox *hcloud.StorageBox, opts hcloud.StorageBoxChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().StorageBox().ChangeProtection(s, storageBox, opts)
	},

	IDOrName: func(storageBox *hcloud.StorageBox) string {
		return fmt.Sprint(storageBox.ID)
	},

	Experimental: experimental.StorageBoxes,
}
