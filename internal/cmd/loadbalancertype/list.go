package loadbalancertype

import (
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Load Balancer Types",
	JSONKeyGetByName:   "load_balancer_types",

	DefaultColumns: []string{"id", "name", "description", "max_services", "max_connections", "max_targets"},

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.LoadBalancerTypeListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		loadBalancerTypes, err := s.LoadBalancerType().AllWithOpts(s, opts)

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

	Schema: func(resources []interface{}) interface{} {
		loadBalancerTypeSchemas := make([]schema.LoadBalancerType, 0, len(resources))
		for _, resource := range resources {
			loadBalancerType := resource.(*hcloud.LoadBalancerType)
			loadBalancerTypeSchemas = append(loadBalancerTypeSchemas, hcloud.SchemaFromLoadBalancerType(loadBalancerType))
		}
		return loadBalancerTypeSchemas
	},
}
