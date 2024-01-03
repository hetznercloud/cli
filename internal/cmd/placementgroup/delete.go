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
	ShortDescription:     "Delete a placement group",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.PlacementGroup().Get(s, idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		placementGroup := resource.(*hcloud.PlacementGroup)
		if _, err := s.PlacementGroup().Delete(s, placementGroup); err != nil {
			return err
		}
		return nil
	},
}
