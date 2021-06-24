package iso

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
	ResourceNamePlural: "isos",
	DefaultColumns:     []string{"id", "name", "description", "type"},

	Fetch: func(ctx context.Context, client hcapi2.Client, listOpts hcloud.ListOpts) ([]interface{}, error) {
		isos, _, err := client.ISO().List(ctx, hcloud.ISOListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, n := range isos {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(_ hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Location{})
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var isoSchemas []schema.ISO
		for _, resource := range resources {
			isoSchemas = append(isoSchemas, util.ISOToSchema(resource.(hcloud.ISO)))
		}
		return isoSchemas
	},
}
