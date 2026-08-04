package placementgroup

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.PlacementGroup]{
	ResourceNameSingular:   "Placement Group",
	ShortDescriptionAdd:    "Add a label to a Placement Group",
	ShortDescriptionRemove: "Remove a label from a Placement Group",
	NameSuggestions:        func(c hcapi2.Client) hcapi2.CompletionFunc { return c.PlacementGroup().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) hcapi2.LabelCompletionFunc { return c.PlacementGroup().LabelKeys },
	Fetch: func(ctx context.Context, s state.State, idOrName string) (*hcloud.PlacementGroup, error) {
		placementGroup, _, err := s.Client().PlacementGroup().Get(ctx, idOrName)
		if err != nil {
			return nil, err
		}
		if placementGroup == nil {
			return nil, fmt.Errorf("Placement Group not found: %s", idOrName)
		}
		return placementGroup, nil
	},
	SetLabels: func(ctx context.Context, s state.State, placementGroup *hcloud.PlacementGroup, labels map[string]string) error {
		opts := hcloud.PlacementGroupUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().PlacementGroup().Update(ctx, placementGroup, opts)
		return err
	},
	GetLabels: func(placementGroup *hcloud.PlacementGroup) map[string]string {
		return placementGroup.Labels
	},
	GetIDOrName: func(placementGroup *hcloud.PlacementGroup) string {
		return strconv.FormatInt(placementGroup.ID, 10)
	},
}
