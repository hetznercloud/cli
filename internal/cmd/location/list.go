package location

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "locations",
	DefaultColumns:     []string{"id", "name", "description", "network_zone", "country", "city"},

	Fetch: func(ctx context.Context, client hcapi2.Client, listOpts hcloud.ListOpts) ([]interface{}, error) {
		locations, _, err := client.Location().List(ctx, hcloud.LocationListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, n := range locations {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(_ hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Location{})
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var locationSchemas []schema.Location
		for _, resource := range resources {
			locationSchemas = append(locationSchemas, util.LocationToSchema(resource.(hcloud.Location)))
		}
		return locationSchemas
	},
}
