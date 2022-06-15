package iso

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
	ResourceNamePlural: "isos",
	DefaultColumns:     []string{"id", "name", "description", "type"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.ISOListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		isos, _, err := client.ISO().List(ctx, opts)

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
			iso := resource.(*hcloud.ISO)
			isoSchemas = append(isoSchemas, util.ISOToSchema(*iso))
		}
		return isoSchemas
	},
}
