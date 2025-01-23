package placementgroup

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "placement group",
	ShortDescriptionAdd:    "Add a label to a placement group",
	ShortDescriptionRemove: "Remove a label from a placement group",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.PlacementGroup().LabelKeys },
	Fetch: func(s state.State, idOrName string) (any, error) {
		placementGroup, _, err := s.Client().PlacementGroup().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if placementGroup == nil {
			return nil, fmt.Errorf("placement group not found: %s", idOrName)
		}
		return placementGroup, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		placementGroup := resource.(*hcloud.PlacementGroup)
		opts := hcloud.PlacementGroupUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().PlacementGroup().Update(s, placementGroup, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		placementGroup := resource.(*hcloud.PlacementGroup)
		return placementGroup.Labels
	},
	GetIDOrName: func(resource any) string {
		placementGroup := resource.(*hcloud.PlacementGroup)
		return strconv.FormatInt(placementGroup.ID, 10)
	},
}
