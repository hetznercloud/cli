package loadbalancer

import (
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.LoadBalancer, schema.LoadBalancer]{
	ResourceNamePlural: "Load Balancer",
	JSONKeyGetByName:   "load_balancers",
	DefaultColumns:     []string{"id", "name", "health", "ipv4", "ipv6", "type", "location", "network_zone", "age"},
	SortOption:         config.OptionSortLoadBalancer,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.LoadBalancer, error) {
		opts := hcloud.LoadBalancerListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().LoadBalancer().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.LoadBalancer], _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.LoadBalancer{}).
			AddFieldFn("ipv4", func(loadBalancer *hcloud.LoadBalancer) string {
				return loadBalancer.PublicNet.IPv4.IP.String()
			}).
			AddFieldFn("ipv6", func(loadBalancer *hcloud.LoadBalancer) string {
				return loadBalancer.PublicNet.IPv6.IP.String()
			}).
			AddFieldFn("type", func(loadBalancer *hcloud.LoadBalancer) string {
				return loadBalancer.LoadBalancerType.Name
			}).
			AddFieldFn("location", func(loadBalancer *hcloud.LoadBalancer) string {
				return loadBalancer.Location.Name
			}).
			AddFieldFn("network_zone", func(loadBalancer *hcloud.LoadBalancer) string {
				return string(loadBalancer.Location.NetworkZone)
			}).
			AddFieldFn("labels", func(loadBalancer *hcloud.LoadBalancer) string {
				return util.LabelsToString(loadBalancer.Labels)
			}).
			AddFieldFn("protection", func(loadBalancer *hcloud.LoadBalancer) string {
				var protection []string
				if loadBalancer.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			}).
			AddFieldFn("created", func(loadBalancer *hcloud.LoadBalancer) string {
				return util.Datetime(loadBalancer.Created)
			}).
			AddFieldFn("age", func(loadBalancer *hcloud.LoadBalancer) string {
				return util.Age(loadBalancer.Created, time.Now())
			}).
			AddFieldFn("health", func(loadBalancer *hcloud.LoadBalancer) string {
				return Health(loadBalancer)
			})
	},

	Schema: hcloud.SchemaFromLoadBalancer,
}

func Health(l *hcloud.LoadBalancer) string {
	healthyCount := 0
	unhealthyCount := 0
	unknownCount := 0

	for _, lbTarget := range l.Targets {
		for _, svcHealth := range lbTarget.HealthStatus {
			switch svcHealth.Status {
			case hcloud.LoadBalancerTargetHealthStatusStatusHealthy:
				healthyCount++

			case hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy:
				unhealthyCount++

			case hcloud.LoadBalancerTargetHealthStatusStatusUnknown:
				unknownCount++
			}
		}
	}

	switch len(l.Targets) * len(l.Services) {
	case healthyCount:
		return string(hcloud.LoadBalancerTargetHealthStatusStatusHealthy)

	case unhealthyCount:
		return string(hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy)

	case unknownCount:
		return string(hcloud.LoadBalancerTargetHealthStatusStatusUnknown)

	default:
		return "mixed"
	}
}
