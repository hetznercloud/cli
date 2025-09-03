package loadbalancer

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DescribeCmd defines a command for describing a LoadBalancer.
var DescribeCmd = base.DescribeCmd[*hcloud.LoadBalancer]{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Describe a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.LoadBalancer, any, error) {
		lb, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return lb, hcloud.SchemaFromLoadBalancer(lb), nil
	},
	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().Bool("expand-targets", false, "Expand all label_selector targets (true, false)")
	},
	PrintText: func(s state.State, cmd *cobra.Command, loadBalancer *hcloud.LoadBalancer, _ base.DescribeWriter) error {
		withLabelSelectorTargets, _ := cmd.Flags().GetBool("expand-targets")
		cmd.Printf("ID:\t\t\t\t%d\n", loadBalancer.ID)
		cmd.Printf("Name:\t\t\t\t%s\n", loadBalancer.Name)
		cmd.Printf("Created:\t\t\t%s (%s)\n", util.Datetime(loadBalancer.Created), humanize.Time(loadBalancer.Created))
		cmd.Printf("Public Net:\n")
		cmd.Printf("  Enabled:\t\t\t%s\n", util.YesNo(loadBalancer.PublicNet.Enabled))
		cmd.Printf("  IPv4:\t\t\t\t%s\n", loadBalancer.PublicNet.IPv4.IP.String())
		cmd.Printf("  IPv4 DNS PTR:\t\t\t%s\n", loadBalancer.PublicNet.IPv4.DNSPtr)
		cmd.Printf("  IPv6:\t\t\t\t%s\n", loadBalancer.PublicNet.IPv6.IP.String())
		cmd.Printf("  IPv6 DNS PTR:\t\t\t%s\n", loadBalancer.PublicNet.IPv6.DNSPtr)

		cmd.Printf("Private Net:\n")
		if len(loadBalancer.PrivateNet) > 0 {
			for _, n := range loadBalancer.PrivateNet {
				cmd.Printf("  - ID:\t\t\t%d\n", n.Network.ID)
				cmd.Printf("    Name:\t\t%s\n", s.Client().Network().Name(n.Network.ID))
				cmd.Printf("    IP:\t\t\t%s\n", n.IP.String())
			}
		} else {
			cmd.Printf("    No Private Network\n")
		}
		cmd.Printf("Algorithm:\t\t\t%s\n", loadBalancer.Algorithm.Type)

		cmd.Printf("Load Balancer Type:\t\t%s (ID: %d)\n", loadBalancer.LoadBalancerType.Name, loadBalancer.LoadBalancerType.ID)
		cmd.Printf("  ID:\t\t\t\t%d\n", loadBalancer.LoadBalancerType.ID)
		cmd.Printf("  Name:\t\t\t\t%s\n", loadBalancer.LoadBalancerType.Name)
		cmd.Printf("  Description:\t\t\t%s\n", loadBalancer.LoadBalancerType.Description)
		cmd.Printf("  Max Services:\t\t\t%d\n", loadBalancer.LoadBalancerType.MaxServices)
		cmd.Printf("  Max Connections:\t\t%d\n", loadBalancer.LoadBalancerType.MaxConnections)
		cmd.Printf("  Max Targets:\t\t\t%d\n", loadBalancer.LoadBalancerType.MaxTargets)
		cmd.Printf("  Max assigned Certificates:\t%d\n", loadBalancer.LoadBalancerType.MaxAssignedCertificates)

		cmd.Printf("Services:\n")
		if len(loadBalancer.Services) == 0 {
			cmd.Print("  No services\n")
		} else {
			for _, service := range loadBalancer.Services {
				cmd.Printf("  - Protocol:\t\t\t%s\n", service.Protocol)
				cmd.Printf("    Listen Port:\t\t%d\n", service.ListenPort)
				cmd.Printf("    Destination Port:\t\t%d\n", service.DestinationPort)
				cmd.Printf("    Proxy Protocol:\t\t%s\n", util.YesNo(service.Proxyprotocol))
				if service.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					cmd.Printf("    Sticky Sessions:\t\t%s\n", util.YesNo(service.HTTP.StickySessions))
					if service.HTTP.StickySessions {
						cmd.Printf("    Sticky Cookie Name:\t\t%s\n", service.HTTP.CookieName)
						cmd.Printf("    Sticky Cookie Lifetime:\t%vs\n", service.HTTP.CookieLifetime.Seconds())
					}
					if service.Protocol == hcloud.LoadBalancerServiceProtocolHTTPS {
						cmd.Printf("    Certificates:\n")
						for _, cert := range service.HTTP.Certificates {
							cmd.Printf("      - ID: \t\t\t%v\n", cert.ID)
						}
					}
				}

				cmd.Printf("    Health Check:\n")
				cmd.Printf("      Protocol:\t\t\t%s\n", service.HealthCheck.Protocol)
				cmd.Printf("      Timeout:\t\t\t%vs\n", service.HealthCheck.Timeout.Seconds())
				cmd.Printf("      Interval:\t\t\tevery %vs\n", service.HealthCheck.Interval.Seconds())
				cmd.Printf("      Retries:\t\t\t%d\n", service.HealthCheck.Retries)
				if service.HealthCheck.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					cmd.Printf("      HTTP Domain:\t\t%s\n", service.HealthCheck.HTTP.Domain)
					cmd.Printf("      HTTP Path:\t\t%s\n", service.HealthCheck.HTTP.Path)
					cmd.Printf("      Response:\t\t%s\n", service.HealthCheck.HTTP.Response)
					cmd.Printf("      TLS:\t\t\t%s\n", util.YesNo(service.HealthCheck.HTTP.TLS))
					cmd.Printf("      Status Codes:\t\t%v\n", service.HealthCheck.HTTP.StatusCodes)
				}
			}
		}

		cmd.Printf("Targets:\n")
		if len(loadBalancer.Targets) == 0 {
			cmd.Print("  No targets\n")
		}
		for _, target := range loadBalancer.Targets {
			cmd.Printf("  - Type:\t\t\t%s\n", target.Type)
			switch target.Type {
			case hcloud.LoadBalancerTargetTypeServer:
				cmd.Printf("    Server:\n")
				cmd.Printf("      ID:\t\t\t%d\n", target.Server.Server.ID)
				cmd.Printf("      Name:\t\t\t%s\n", s.Client().Server().ServerName(target.Server.Server.ID))
				cmd.Printf("    Use Private IP:\t\t%s\n", util.YesNo(target.UsePrivateIP))
				cmd.Printf("    Status:\n")
				for _, healthStatus := range target.HealthStatus {
					cmd.Printf("    - Service:\t\t\t%d\n", healthStatus.ListenPort)
					cmd.Printf("      Status:\t\t\t%s\n", healthStatus.Status)
				}
			case hcloud.LoadBalancerTargetTypeLabelSelector:
				cmd.Printf("    Label Selector:\t\t%s\n", target.LabelSelector.Selector)
				cmd.Printf("      Targets: (%d)\n", len(target.Targets))
				if len(target.Targets) == 0 {
					cmd.Print("      No targets\n")
				}
				if !withLabelSelectorTargets {
					continue
				}
				for _, lbtarget := range target.Targets {
					cmd.Printf("      - Type:\t\t\t\t%s\n", lbtarget.Type)
					cmd.Printf("        Server ID:\t\t\t%d\n", lbtarget.Server.Server.ID)
					cmd.Printf("        Status:\n")
					for _, healthStatus := range lbtarget.HealthStatus {
						cmd.Printf("          - Service:\t\t\t%d\n", healthStatus.ListenPort)
						cmd.Printf("            Status:\t\t\t%s\n", healthStatus.Status)
					}
				}
			case hcloud.LoadBalancerTargetTypeIP:
				cmd.Printf("    IP:\t\t\t\t%s\n", target.IP.IP)
				cmd.Printf("    Status:\n")
				for _, healthStatus := range target.HealthStatus {
					cmd.Printf("    - Service:\t\t\t%d\n", healthStatus.ListenPort)
					cmd.Printf("      Status:\t\t\t%s\n", healthStatus.Status)
				}
			}
		}

		cmd.Printf("Traffic:\n")
		cmd.Printf("  Outgoing:\t%v\n", humanize.IBytes(loadBalancer.OutgoingTraffic))
		cmd.Printf("  Ingoing:\t%v\n", humanize.IBytes(loadBalancer.IngoingTraffic))
		cmd.Printf("  Included:\t%v\n", humanize.IBytes(loadBalancer.IncludedTraffic))

		cmd.Printf("Protection:\n")
		cmd.Printf("  Delete:\t%s\n", util.YesNo(loadBalancer.Protection.Delete))

		cmd.Print("Labels:\n")
		if len(loadBalancer.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range util.IterateInOrder(loadBalancer.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
