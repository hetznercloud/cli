package placementgroup

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "placement group",
	ResourceNamePlural:   "placement groups",
	ShortDescription:     "Delete a placement group",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().PlacementGroup().Get(s, idOrName)
	},
	Delete: func(s state.State, _ *cobra.Command, resource interface{}) (*hcloud.Action, error) {
		placementGroup := resource.(*hcloud.PlacementGroup)
		_, err := s.Client().PlacementGroup().Delete(s, placementGroup)
		return nil, err
	},
}
