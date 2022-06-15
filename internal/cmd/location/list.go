package location

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "locations",
	DefaultColumns:     []string{"id", "name", "description", "network_zone", "country", "city"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.LocationListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		locations, _, err := client.Location().List(ctx, opts)

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
			location := resource.(*hcloud.Location)
			locationSchemas = append(locationSchemas, util.LocationToSchema(*location))
		}
		return locationSchemas
	},
}
