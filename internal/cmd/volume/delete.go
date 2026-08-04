package volume

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.Volume]{
	ResourceNameSingular: "Volume",
	ResourceNamePlural:   "Volumes",
	ShortDescription:     "Delete a Volume",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Volume().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Volume, *hcloud.Response, error) {
		return s.Client().Volume().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, volume *hcloud.Volume) ([]*hcloud.Action, error) {
		_, err := s.Client().Volume().Delete(cmd.Context(), volume)
		return nil, err
	},
}
