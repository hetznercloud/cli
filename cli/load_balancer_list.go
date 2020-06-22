package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var loadBalancerListTableOutput *tableOutput

func init() {
	loadBalancerListTableOutput = describeLoadBalancerListTableOutput(nil)
}

func newLoadBalancerListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Load Balancers",
		Long: listLongDescription(
			"Displays a list of Load Balancers.",
			loadBalancerListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(loadBalancerListTableOutput.Columns()), outputOptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runLoadBalancerList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.LoadBalancerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	loadBalancers, err := cli.Client().LoadBalancer.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}
	if outOpts.IsSet("json") {
		var loadBalancerSchemas []schema.LoadBalancer
		for _, loadBalancer := range loadBalancers {
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
				LoadBalancerType: loadBalancerTypeToSchema(*loadBalancer.LoadBalancerType),
				Location:         locationToSchema(*loadBalancer.Location),
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
		return describeJSON(loadBalancerSchemas)
	}
	cols := []string{"id", "name", "ipv4", "ipv6", "type", "location", "network_zone"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := describeLoadBalancerListTableOutput(cli)
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, loadBalancer := range loadBalancers {
		tw.Write(cols, loadBalancer)
	}
	tw.Flush()
	return nil
}

func describeLoadBalancerListTableOutput(cli *CLI) *tableOutput {
	return newTableOutput().
		AddAllowedFields(hcloud.LoadBalancer{}).
		AddFieldOutputFn("ipv4", fieldOutputFn(func(obj interface{}) string {
			loadbalancer := obj.(*hcloud.LoadBalancer)
			return loadbalancer.PublicNet.IPv4.IP.String()
		})).
		AddFieldOutputFn("ipv6", fieldOutputFn(func(obj interface{}) string {
			loadbalancer := obj.(*hcloud.LoadBalancer)
			return loadbalancer.PublicNet.IPv6.IP.String()
		})).
		AddFieldOutputFn("type", fieldOutputFn(func(obj interface{}) string {
			loadbalancer := obj.(*hcloud.LoadBalancer)
			return loadbalancer.LoadBalancerType.Name
		})).
		AddFieldOutputFn("location", fieldOutputFn(func(obj interface{}) string {
			loadbalancer := obj.(*hcloud.LoadBalancer)
			return loadbalancer.Location.Name
		})).
		AddFieldOutputFn("network_zone", fieldOutputFn(func(obj interface{}) string {
			loadbalancer := obj.(*hcloud.LoadBalancer)
			return string(loadbalancer.Location.NetworkZone)
		})).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			loadBalancer := obj.(*hcloud.LoadBalancer)
			return labelsToString(loadBalancer.Labels)
		})).
		AddFieldOutputFn("protection", fieldOutputFn(func(obj interface{}) string {
			loadBalancer := obj.(*hcloud.LoadBalancer)
			var protection []string
			if loadBalancer.Protection.Delete {
				protection = append(protection, "delete")
			}
			return strings.Join(protection, ", ")
		})).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			loadBalancer := obj.(*hcloud.LoadBalancer)
			return datetime(loadBalancer.Created)
		}))
}
