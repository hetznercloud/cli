package datacenter

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "datacenters",
	DefaultColumns:     []string{"id", "name", "description", "location"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		datacenters, _, err := client.Datacenter().List(ctx, hcloud.DatacenterListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, n := range datacenters {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(_ hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Datacenter{}).
			AddFieldFn("location", output.FieldFn(func(obj interface{}) string {
				datacenter := obj.(*hcloud.Datacenter)
				return datacenter.Location.Name
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var certSchemas []schema.Datacenter
		for _, resource := range resources {
			cert := resource.(*hcloud.Datacenter)
			certSchemas = append(certSchemas, util.DatacenterToSchema(*cert))
		}

		return util.DescribeJSON(certSchemas)
	},
}
