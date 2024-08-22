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

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Load Balancer",
	JSONKeyGetByName:   "load_balancers",
	DefaultColumns:     []string{"id", "name", "health", "ipv4", "ipv6", "type", "location", "network_zone", "age"},
	SortOption:         config.OptionSortLoadBalancer,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.LoadBalancerListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		loadBalancers, err := s.Client().LoadBalancer().AllWithOpts(s, opts)

		var resources []interface{}
		for _, r := range loadBalancers {
			resources = append(resources, r)
		}
		return resources, err
	},

	OutputTable: func(hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.LoadBalancer{}).
			AddFieldFn("ipv4", output.FieldFn(func(obj interface{}) string {
				loadbalancer := obj.(*hcloud.LoadBalancer)
				return loadbalancer.PublicNet.IPv4.IP.String()
			})).
			AddFieldFn("ipv6", output.FieldFn(func(obj interface{}) string {
				loadbalancer := obj.(*hcloud.LoadBalancer)
				return loadbalancer.PublicNet.IPv6.IP.String()
			})).
			AddFieldFn("type", output.FieldFn(func(obj interface{}) string {
				loadbalancer := obj.(*hcloud.LoadBalancer)
				return loadbalancer.LoadBalancerType.Name
			})).
			AddFieldFn("location", output.FieldFn(func(obj interface{}) string {
				loadbalancer := obj.(*hcloud.LoadBalancer)
				return loadbalancer.Location.Name
			})).
			AddFieldFn("network_zone", output.FieldFn(func(obj interface{}) string {
				loadbalancer := obj.(*hcloud.LoadBalancer)
				return string(loadbalancer.Location.NetworkZone)
			})).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				loadBalancer := obj.(*hcloud.LoadBalancer)
				return util.LabelsToString(loadBalancer.Labels)
			})).
			AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
				loadBalancer := obj.(*hcloud.LoadBalancer)
				var protection []string
				if loadBalancer.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				loadBalancer := obj.(*hcloud.LoadBalancer)
				return util.Datetime(loadBalancer.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				loadBalancer := obj.(*hcloud.LoadBalancer)
				return util.Age(loadBalancer.Created, time.Now())
			})).
			AddFieldFn("health", output.FieldFn(func(obj interface{}) string {
				loadBalancer := obj.(*hcloud.LoadBalancer)
				return Health(loadBalancer)
			}))
	},

	Schema: func(resources []interface{}) interface{} {
		loadBalancerSchemas := make([]schema.LoadBalancer, 0, len(resources))
		for _, resource := range resources {
			loadBalancer := resource.(*hcloud.LoadBalancer)
			loadBalancerSchemas = append(loadBalancerSchemas, hcloud.SchemaFromLoadBalancer(loadBalancer))
		}
		return loadBalancerSchemas
	},
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

			default:
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
