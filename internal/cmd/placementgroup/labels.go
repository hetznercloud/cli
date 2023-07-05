package placementgroup

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "placement group",
	ShortDescriptionAdd:    "Add a label to a placement group",
	ShortDescriptionRemove: "Remove a label from a placement group",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.PlacementGroup().LabelKeys },
	FetchLabels: func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int64, error) {
		placementGroup, _, err := client.PlacementGroup().Get(ctx, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if placementGroup == nil {
			return nil, 0, fmt.Errorf("placement group not found: %s", idOrName)
		}
		return placementGroup.Labels, placementGroup.ID, nil
	},
	SetLabels: func(ctx context.Context, client hcapi2.Client, id int64, labels map[string]string) error {
		opts := hcloud.PlacementGroupUpdateOpts{
			Labels: labels,
		}
		_, _, err := client.PlacementGroup().Update(ctx, &hcloud.PlacementGroup{ID: id}, opts)
		return err
	},
}
