package placementgroup

import (
	"context"

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
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.PlacementGroup().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, _ state.ActionWaiter, cmd *cobra.Command, resource interface{}) error {
		placementGroup := resource.(*hcloud.PlacementGroup)
		if _, err := client.PlacementGroup().Delete(ctx, placementGroup); err != nil {
			return err
		}
		return nil
	},
}
