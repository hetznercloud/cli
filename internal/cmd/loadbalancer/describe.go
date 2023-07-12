package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

// DescribeCmd defines a command for describing a LoadBalancer.
var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Describe a Load Balancer",
	JSONKeyGetByID:       "load_balancer",
	JSONKeyGetByName:     "load_balancers",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.LoadBalancer().Get(ctx, idOrName)
	},
	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().Bool("expand-targets", false, "Expand all label_selector targets")
	},
	PrintText: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		withLabelSelectorTargets, _ := cmd.Flags().GetBool("expand-targets")
		loadBalancer := resource.(*hcloud.LoadBalancer)
		fmt.Printf("ID:\t\t\t\t%d\n", loadBalancer.ID)
		fmt.Printf("Name:\t\t\t\t%s\n", loadBalancer.Name)
		fmt.Printf("Created:\t\t\t%s (%s)\n", util.Datetime(loadBalancer.Created), humanize.Time(loadBalancer.Created))
		fmt.Printf("Public Net:\n")
		fmt.Printf("  Enabled:\t\t\t%s\n", util.YesNo(loadBalancer.PublicNet.Enabled))
		fmt.Printf("  IPv4:\t\t\t\t%s\n", loadBalancer.PublicNet.IPv4.IP.String())
		fmt.Printf("  IPv4 DNS PTR:\t\t\t%s\n", loadBalancer.PublicNet.IPv4.DNSPtr)
		fmt.Printf("  IPv6:\t\t\t\t%s\n", loadBalancer.PublicNet.IPv6.IP.String())
		fmt.Printf("  IPv6 DNS PTR:\t\t\t%s\n", loadBalancer.PublicNet.IPv6.DNSPtr)

		fmt.Printf("Private Net:\n")
		if len(loadBalancer.PrivateNet) > 0 {
			for _, n := range loadBalancer.PrivateNet {
				fmt.Printf("  - ID:\t\t\t%d\n", n.Network.ID)
				fmt.Printf("    Name:\t\t%s\n", client.Network().Name(n.Network.ID))
				fmt.Printf("    IP:\t\t\t%s\n", n.IP.String())
			}
		} else {
			fmt.Printf("    No Private Network\n")
		}
		fmt.Printf("Algorithm:\t\t\t%s\n", loadBalancer.Algorithm.Type)

		fmt.Printf("Load Balancer Type:\t\t%s (ID: %d)\n", loadBalancer.LoadBalancerType.Name, loadBalancer.LoadBalancerType.ID)
		fmt.Printf("  ID:\t\t\t\t%d\n", loadBalancer.LoadBalancerType.ID)
		fmt.Printf("  Name:\t\t\t\t%s\n", loadBalancer.LoadBalancerType.Name)
		fmt.Printf("  Description:\t\t\t%s\n", loadBalancer.LoadBalancerType.Description)
		fmt.Printf("  Max Services:\t\t\t%d\n", loadBalancer.LoadBalancerType.MaxServices)
		fmt.Printf("  Max Connections:\t\t%d\n", loadBalancer.LoadBalancerType.MaxConnections)
		fmt.Printf("  Max Targets:\t\t\t%d\n", loadBalancer.LoadBalancerType.MaxTargets)
		fmt.Printf("  Max assigned Certificates:\t%d\n", loadBalancer.LoadBalancerType.MaxAssignedCertificates)

		fmt.Printf("Services:\n")
		if len(loadBalancer.Services) == 0 {
			fmt.Print("  No services\n")
		} else {
			for _, service := range loadBalancer.Services {
				fmt.Printf("  - Protocol:\t\t\t%s\n", service.Protocol)
				fmt.Printf("    Listen Port:\t\t%d\n", service.ListenPort)
				fmt.Printf("    Destination Port:\t\t%d\n", service.DestinationPort)
				fmt.Printf("    Proxy Protocol:\t\t%s\n", util.YesNo(service.Proxyprotocol))
				if service.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					fmt.Printf("    Sticky Sessions:\t\t%s\n", util.YesNo(service.HTTP.StickySessions))
					if service.HTTP.StickySessions {
						fmt.Printf("    Sticky Cookie Name:\t\t%s\n", service.HTTP.CookieName)
						fmt.Printf("    Sticky Cookie Lifetime:\t%vs\n", service.HTTP.CookieLifetime.Seconds())
					}
					if service.Protocol == hcloud.LoadBalancerServiceProtocolHTTPS {
						fmt.Printf("    Certificates:\n")
						for _, cert := range service.HTTP.Certificates {
							fmt.Printf("      - ID: \t\t\t%v\n", cert.ID)
						}
					}
				}

				fmt.Printf("    Health Check:\n")
				fmt.Printf("      Protocol:\t\t\t%s\n", service.HealthCheck.Protocol)
				fmt.Printf("      Timeout:\t\t\t%vs\n", service.HealthCheck.Timeout.Seconds())
				fmt.Printf("      Interval:\t\t\tevery %vs\n", service.HealthCheck.Interval.Seconds())
				fmt.Printf("      Retries:\t\t\t%d\n", service.HealthCheck.Retries)
				if service.HealthCheck.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					fmt.Printf("      HTTP Domain:\t\t%s\n", service.HealthCheck.HTTP.Domain)
					fmt.Printf("      HTTP Path:\t\t%s\n", service.HealthCheck.HTTP.Path)
					fmt.Printf("      Response:\t\t%s\n", service.HealthCheck.HTTP.Response)
					fmt.Printf("      TLS:\t\t\t%s\n", util.YesNo(service.HealthCheck.HTTP.TLS))
					fmt.Printf("      Status Codes:\t\t%v\n", service.HealthCheck.HTTP.StatusCodes)
				}
			}
		}

		fmt.Printf("Targets:\n")
		if len(loadBalancer.Targets) == 0 {
			fmt.Print("  No targets\n")
		}
		for _, target := range loadBalancer.Targets {
			fmt.Printf("  - Type:\t\t\t%s\n", target.Type)
			switch target.Type {
			case hcloud.LoadBalancerTargetTypeServer:
				fmt.Printf("    Server:\n")
				fmt.Printf("      ID:\t\t\t%d\n", target.Server.Server.ID)
				fmt.Printf("      Name:\t\t\t%s\n", client.Server().ServerName(target.Server.Server.ID))
				fmt.Printf("    Use Private IP:\t\t%s\n", util.YesNo(target.UsePrivateIP))
				fmt.Printf("    Status:\n")
				for _, healthStatus := range target.HealthStatus {
					fmt.Printf("    - Service:\t\t\t%d\n", healthStatus.ListenPort)
					fmt.Printf("      Status:\t\t\t%s\n", healthStatus.Status)
				}
			case hcloud.LoadBalancerTargetTypeLabelSelector:
				fmt.Printf("    Label Selector:\t\t%s\n", target.LabelSelector.Selector)
				fmt.Printf("      Targets: (%d)\n", len(target.Targets))
				if len(target.Targets) == 0 {
					fmt.Print("      No targets\n")
				}
				if !withLabelSelectorTargets {
					continue
				}
				for _, lbtarget := range target.Targets {
					fmt.Printf("      - Type:\t\t\t\t%s\n", lbtarget.Type)
					fmt.Printf("        Server ID:\t\t\t%d\n", lbtarget.Server.Server.ID)
					fmt.Printf("        Status:\n")
					for _, healthStatus := range lbtarget.HealthStatus {
						fmt.Printf("          - Service:\t\t\t%d\n", healthStatus.ListenPort)
						fmt.Printf("            Status:\t\t\t%s\n", healthStatus.Status)
					}
				}
			case hcloud.LoadBalancerTargetTypeIP:
				fmt.Printf("    IP:\t\t\t\t%s\n", target.IP.IP)
				fmt.Printf("    Status:\n")
				for _, healthStatus := range target.HealthStatus {
					fmt.Printf("    - Service:\t\t\t%d\n", healthStatus.ListenPort)
					fmt.Printf("      Status:\t\t\t%s\n", healthStatus.Status)
				}
			}
		}

		fmt.Printf("Traffic:\n")
		fmt.Printf("  Outgoing:\t%v\n", humanize.IBytes(loadBalancer.OutgoingTraffic))
		fmt.Printf("  Ingoing:\t%v\n", humanize.IBytes(loadBalancer.IngoingTraffic))
		fmt.Printf("  Included:\t%v\n", humanize.IBytes(loadBalancer.IncludedTraffic))

		fmt.Printf("Protection:\n")
		fmt.Printf("  Delete:\t%s\n", util.YesNo(loadBalancer.Protection.Delete))

		fmt.Print("Labels:\n")
		if len(loadBalancer.Labels) == 0 {
			fmt.Print("  No labels\n")
		} else {
			for key, value := range loadBalancer.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
