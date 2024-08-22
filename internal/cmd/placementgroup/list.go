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

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Placement Groups",
	JSONKeyGetByName:   "placement_groups",
	DefaultColumns:     []string{"id", "name", "servers", "type", "age"},
	SortOption:         config.OptionSortPlacementGroup,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.PlacementGroupListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		placementGroups, err := s.Client().PlacementGroup().AllWithOpts(s, opts)

		var resources []interface{}
		for _, n := range placementGroups {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.PlacementGroup{}).
			AddFieldFn("servers", output.FieldFn(func(obj interface{}) string {
				placementGroup := obj.(*hcloud.PlacementGroup)
				count := len(placementGroup.Servers)
				if count == 1 {
					return fmt.Sprintf("%d server", count)
				}
				return fmt.Sprintf("%d servers", count)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				placementGroup := obj.(*hcloud.PlacementGroup)
				return util.Datetime(placementGroup.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				placementGroup := obj.(*hcloud.PlacementGroup)
				return util.Age(placementGroup.Created, time.Now())
			}))
	},

	Schema: func(resources []interface{}) interface{} {
		placementGroupSchemas := make([]schema.PlacementGroup, 0, len(resources))
		for _, resource := range resources {
			placementGroup := resource.(*hcloud.PlacementGroup)
			placementGroupSchemas = append(placementGroupSchemas, hcloud.SchemaFromPlacementGroup(placementGroup))
		}
		return placementGroupSchemas
	},
}
