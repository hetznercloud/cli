package placementgroup

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "placement groups",
	DefaultColumns:     []string{"id", "name", "servers_count", "type"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		placementGroups, _, err := client.PlacementGroup().List(ctx, hcloud.PlacementGroupListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, n := range placementGroups {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.PlacementGroup{}).
			AddFieldFn("servers_count", output.FieldFn(func(obj interface{}) string {
				placementGroup := obj.(*hcloud.PlacementGroup)
				count := len(placementGroup.Servers)
				if count == 1 {
					return fmt.Sprintf("%d Server", count)
				}
				return fmt.Sprintf("%d Servers", count)
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var placementGroupSchemas []schema.PlacementGroup
		for _, resource := range resources {
			placementGroup := resource.(*hcloud.PlacementGroup)
			placementGroupSchema := schema.PlacementGroup{
				ID:      placementGroup.ID,
				Name:    placementGroup.Name,
				Labels:  placementGroup.Labels,
				Created: placementGroup.Created,
				Servers: placementGroup.Servers,
				Type:    string(placementGroup.Type),
			}

			placementGroupSchemas = append(placementGroupSchemas, placementGroupSchema)
		}
		return placementGroupSchemas
	},
}
