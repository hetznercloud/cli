package loadbalancertype

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Load Balancer Types",

	DefaultColumns: []string{"id", "name", "description", "max_services", "max_connections", "max_targets"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		loadBalancerTypes, _, err := client.LoadBalancerType().List(ctx, hcloud.LoadBalancerTypeListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, r := range loadBalancerTypes {
			resources = append(resources, r)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.LoadBalancerType{})
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var loadBalancerTypeSchemas []schema.LoadBalancerType
		for _, resource := range resources {
			loadBalancerType := resource.(*hcloud.LoadBalancerType)
			loadBalancerTypeSchemas = append(loadBalancerTypeSchemas, util.LoadBalancerTypeToSchema(*loadBalancerType))
		}
		return loadBalancerTypeSchemas
	},
}
