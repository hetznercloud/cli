package placementgroup

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "placement group",
	ShortDescription:     "Update a placement group",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.PlacementGroup().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Placement group name")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		placementGroup := resource.(*hcloud.PlacementGroup)
		updOpts := hcloud.PlacementGroupUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := client.PlacementGroup().Update(ctx, placementGroup, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
