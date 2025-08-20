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

var ListCmd = base.ListCmd[*hcloud.LoadBalancerType, schema.LoadBalancerType]{
	ResourceNamePlural: "Load Balancer Types",
	JSONKeyGetByName:   "load_balancer_types",
	DefaultColumns:     []string{"id", "name", "description", "max_services", "max_connections", "max_targets"},
	SortOption:         nil, // Load Balancer Types do not support sorting

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.LoadBalancerType, error) {
		opts := hcloud.LoadBalancerTypeListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().LoadBalancerType().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.LoadBalancerType{})
	},

	Schema: hcloud.SchemaFromLoadBalancerType,
}
