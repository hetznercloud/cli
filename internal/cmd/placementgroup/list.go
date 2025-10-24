package placementgroup

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.PlacementGroup, schema.PlacementGroup]{
	ResourceNamePlural: "Placement Groups",
	JSONKeyGetByName:   "placement_groups",
	DefaultColumns:     []string{"id", "name", "servers", "type", "age"},
	SortOption:         config.OptionSortPlacementGroup,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.PlacementGroup, error) {
		opts := hcloud.PlacementGroupListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().PlacementGroup().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.PlacementGroup], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.PlacementGroup{}).
			AddFieldFn("servers", func(placementGroup *hcloud.PlacementGroup) string {
				count := len(placementGroup.Servers)
				if count == 1 {
					return fmt.Sprintf("%d server", count)
				}
				return fmt.Sprintf("%d servers", count)
			}).
			AddFieldFn("created", func(placementGroup *hcloud.PlacementGroup) string {
				return util.Datetime(placementGroup.Created)
			}).
			AddFieldFn("age", func(placementGroup *hcloud.PlacementGroup) string {
				return util.Age(placementGroup.Created, time.Now())
			})
	},

	Schema: hcloud.SchemaFromPlacementGroup,
}
