package storagebox

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.StorageBox]{
	ResourceNameSingular: "Storage Box",
	ResourceNamePlural:   "Storage Boxes",
	ShortDescription:     "Delete a Storage Box",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.StorageBox().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.StorageBox, *hcloud.Response, error) {
		return s.Client().StorageBox().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, storageBox *hcloud.StorageBox) ([]*hcloud.Action, error) {
		result, _, err := s.Client().StorageBox().Delete(cmd.Context(), storageBox)
		return []*hcloud.Action{result.Action}, err
	},
}
