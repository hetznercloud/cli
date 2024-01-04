package servertype

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Server Types",
	JSONKeyGetByName:   "server_types",

	DefaultColumns: []string{"id", "name", "cores", "cpu_type", "architecture", "memory", "disk", "storage_type", "traffic"},

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.ServerTypeListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		servers, err := s.Client().ServerType().AllWithOpts(s, opts)

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

	Schema: func(resources []interface{}) interface{} {
		serverTypeSchemas := make([]schema.ServerType, 0, len(resources))
		for _, resource := range resources {
			serverType := resource.(*hcloud.ServerType)
			serverTypeSchemas = append(serverTypeSchemas, hcloud.SchemaFromServerType(serverType))
		}
		return serverTypeSchemas
	},
}
