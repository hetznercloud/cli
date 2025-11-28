package placementgroup

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.PlacementGroup]{
	ResourceNameSingular: "Placement Group",
	ResourceNamePlural:   "Placement Groups",
	ShortDescription:     "Delete a Placement Group",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.PlacementGroup, *hcloud.Response, error) {
		return s.Client().PlacementGroup().Get(s, idOrName)
	},
	Delete: func(s state.State, _ *cobra.Command, placementGroup *hcloud.PlacementGroup) (*hcloud.Action, error) {
		_, err := s.Client().PlacementGroup().Delete(s, placementGroup)
		return nil, err
	},
}
