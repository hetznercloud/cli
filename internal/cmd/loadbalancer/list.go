package loadbalancer

import (
	"context"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Load Balancer",

	DefaultColumns: []string{"id", "name", "ipv4", "ipv6", "type", "location", "network_zone"},
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		loadBalancers, _, err := client.LoadBalancer().List(ctx, hcloud.LoadBalancerListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, r := range loadBalancers {
			resources = append(resources, r)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
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
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var loadBalancerSchemas []schema.LoadBalancer
		for _, resource := range resources {
			loadBalancer := resource.(*hcloud.LoadBalancer)
			loadBalancerSchema := schema.LoadBalancer{
				ID:   loadBalancer.ID,
				Name: loadBalancer.Name,
				PublicNet: schema.LoadBalancerPublicNet{
					Enabled: loadBalancer.PublicNet.Enabled,
					IPv4: schema.LoadBalancerPublicNetIPv4{
						IP: loadBalancer.PublicNet.IPv4.IP.String(),
					},
					IPv6: schema.LoadBalancerPublicNetIPv6{
						IP: loadBalancer.PublicNet.IPv6.IP.String(),
					},
				},
				Created:          loadBalancer.Created,
				Labels:           loadBalancer.Labels,
				LoadBalancerType: util.LoadBalancerTypeToSchema(*loadBalancer.LoadBalancerType),
				Location:         util.LocationToSchema(*loadBalancer.Location),
				IncludedTraffic:  loadBalancer.IncludedTraffic,
				OutgoingTraffic:  &loadBalancer.OutgoingTraffic,
				IngoingTraffic:   &loadBalancer.IngoingTraffic,
				Protection: schema.LoadBalancerProtection{
					Delete: loadBalancer.Protection.Delete,
				},
				Algorithm: schema.LoadBalancerAlgorithm{Type: string(loadBalancer.Algorithm.Type)},
			}
			for _, service := range loadBalancer.Services {
				serviceSchema := schema.LoadBalancerService{
					Protocol:        string(service.Protocol),
					ListenPort:      service.ListenPort,
					DestinationPort: service.DestinationPort,
					Proxyprotocol:   service.Proxyprotocol,
					HealthCheck: &schema.LoadBalancerServiceHealthCheck{
						Protocol: string(service.HealthCheck.Protocol),
						Port:     service.HealthCheck.Port,
						Interval: int(service.HealthCheck.Interval.Seconds()),
						Timeout:  int(service.HealthCheck.Timeout.Seconds()),
						Retries:  service.HealthCheck.Retries,
					},
				}
				if service.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					serviceSchema.HTTP = &schema.LoadBalancerServiceHTTP{
						StickySessions: service.HTTP.StickySessions,
						CookieName:     service.HTTP.CookieName,
						CookieLifetime: int(service.HTTP.CookieLifetime.Seconds()),
						RedirectHTTP:   service.HTTP.RedirectHTTP,
					}
				}
				if service.HealthCheck.HTTP != nil {
					serviceSchema.HealthCheck.HTTP = &schema.LoadBalancerServiceHealthCheckHTTP{
						Domain:      service.HealthCheck.HTTP.Domain,
						Path:        service.HealthCheck.HTTP.Path,
						StatusCodes: service.HealthCheck.HTTP.StatusCodes,
						TLS:         service.HealthCheck.HTTP.TLS,
						Response:    service.HealthCheck.HTTP.Response,
					}
				}
				loadBalancerSchema.Services = append(loadBalancerSchema.Services, serviceSchema)
			}
			for _, target := range loadBalancer.Targets {
				targetSchema := schema.LoadBalancerTarget{
					Type:         string(target.Type),
					UsePrivateIP: target.UsePrivateIP,
				}
				if target.Type == hcloud.LoadBalancerTargetTypeServer {
					targetSchema.Server = &schema.LoadBalancerTargetServer{ID: target.Server.Server.ID}
				}
				if target.Type == hcloud.LoadBalancerTargetTypeLabelSelector {
					targetSchema.LabelSelector = &schema.LoadBalancerTargetLabelSelector{Selector: target.LabelSelector.Selector}
				}
				if target.Type == hcloud.LoadBalancerTargetTypeIP {
					targetSchema.IP = &schema.LoadBalancerTargetIP{IP: target.IP.IP}
				}
				for _, healthStatus := range target.HealthStatus {
					targetSchema.HealthStatus = append(targetSchema.HealthStatus, schema.LoadBalancerTargetHealthStatus{
						ListenPort: healthStatus.ListenPort,
						Status:     string(healthStatus.Status),
					})
				}
				loadBalancerSchema.Targets = append(loadBalancerSchema.Targets, targetSchema)
			}

			loadBalancerSchemas = append(loadBalancerSchemas, loadBalancerSchema)
		}
		return loadBalancerSchemas
	},
}
