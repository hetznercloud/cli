package servertype

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Server Types",

	DefaultColumns: []string{"id", "name", "cores", "cpu_type", "memory", "disk", "storage_type"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		servers, _, err := client.ServerType().List(ctx, hcloud.ServerTypeListOpts{ListOpts: listOpts})

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
			}))
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
