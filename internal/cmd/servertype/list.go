package servertype

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Server Types",

	DefaultColumns: []string{"id", "name", "cores", "cpu_type", "architecture", "memory", "disk", "storage_type", "traffic"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.ServerTypeListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		servers, _, err := client.ServerType().List(ctx, opts)

		var resources []interface{}
		for _, r := range servers {
			resources = append(resources, r)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.ServerType{}).
			AddFieldAlias("storagetype", "storage type").
			AddFieldFn("memory", output.FieldFn(func(obj interface{}) string {
				serverType := obj.(*hcloud.ServerType)
				return fmt.Sprintf("%.1f GB", serverType.Memory)
			})).
			AddFieldFn("disk", output.FieldFn(func(obj interface{}) string {
				serverType := obj.(*hcloud.ServerType)
				return fmt.Sprintf("%d GB", serverType.Disk)
			})).
			AddFieldFn("traffic", func(obj interface{}) string {
				serverType := obj.(*hcloud.ServerType)
				return fmt.Sprintf("%d TB", serverType.IncludedTraffic/util.Tebibyte)
			})
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var serverTypeSchemas []schema.ServerType
		for _, resource := range resources {
			serverType := resource.(*hcloud.ServerType)
			serverTypeSchemas = append(serverTypeSchemas, util.ServerTypeToSchema(*serverType))
		}
		return serverTypeSchemas
	},
}
