package loadbalancer

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/loadbalancertype"
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
	PrintText: func(s state.State, cmd *cobra.Command, out io.Writer, loadBalancer *hcloud.LoadBalancer) error {
		withLabelSelectorTargets, _ := cmd.Flags().GetBool("expand-targets")

		fmt.Fprintf(out, "ID:\t%d\n", loadBalancer.ID)
		fmt.Fprintf(out, "Name:\t%s\n", loadBalancer.Name)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(loadBalancer.Created), humanize.Time(loadBalancer.Created))
		fmt.Fprintf(out, "Public Net:\t\n")
		fmt.Fprintf(out, "  Enabled:\t%s\n", util.YesNo(loadBalancer.PublicNet.Enabled))
		fmt.Fprintf(out, "  IPv4:\t%s\n", loadBalancer.PublicNet.IPv4.IP.String())
		fmt.Fprintf(out, "  IPv4 DNS PTR:\t%s\n", loadBalancer.PublicNet.IPv4.DNSPtr)
		fmt.Fprintf(out, "  IPv6:\t%s\n", loadBalancer.PublicNet.IPv6.IP.String())
		fmt.Fprintf(out, "  IPv6 DNS PTR:\t%s\n", loadBalancer.PublicNet.IPv6.DNSPtr)

		if len(loadBalancer.PrivateNet) > 0 {
			fmt.Fprintf(out, "Private Net:\t\n")
			for _, n := range loadBalancer.PrivateNet {
				fmt.Fprintf(out, "  - ID:\t%d\n", n.Network.ID)
				fmt.Fprintf(out, "    Name:\t%s\n", s.Client().Network().Name(n.Network.ID))
				fmt.Fprintf(out, "    IP:\t%s\n", n.IP.String())
			}
		} else {
			fmt.Fprintf(out, "Private Net:\tNo Private Network\n")
		}
		fmt.Fprintf(out, "Algorithm:\t%s\n", loadBalancer.Algorithm.Type)

		fmt.Fprintf(out, "Load Balancer Type:\t\n")
		fmt.Fprintf(out, "%s", util.PrefixLines(loadbalancertype.DescribeLoadBalancerType(s, loadBalancer.LoadBalancerType, true), "  "))

		if len(loadBalancer.Services) == 0 {
			fmt.Fprintf(out, "Services:\tNo services\n")
		} else {
			fmt.Fprintf(out, "Services:\t\n")
			for _, service := range loadBalancer.Services {
				fmt.Fprintf(out, "  - Protocol:\t%s\n", service.Protocol)
				fmt.Fprintf(out, "    Listen Port:\t%d\n", service.ListenPort)
				fmt.Fprintf(out, "    Destination Port:\t%d\n", service.DestinationPort)
				fmt.Fprintf(out, "    Proxy Protocol:\t%s\n", util.YesNo(service.Proxyprotocol))
				if service.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					fmt.Fprintf(out, "    Sticky Sessions:\t%s\n", util.YesNo(service.HTTP.StickySessions))
					if service.HTTP.StickySessions {
						fmt.Fprintf(out, "    Sticky Cookie Name:\t%s\n", service.HTTP.CookieName)
						fmt.Fprintf(out, "    Sticky Cookie Lifetime:\t%vs\n", service.HTTP.CookieLifetime.Seconds())
					}
					if service.Protocol == hcloud.LoadBalancerServiceProtocolHTTPS {
						fmt.Fprintf(out, "    Certificates:\n")
						for _, cert := range service.HTTP.Certificates {
							fmt.Fprintf(out, "      - ID:\t%v\n", cert.ID)
						}
					}
				}

				fmt.Fprintf(out, "    Health Check:\n")
				fmt.Fprintf(out, "      Protocol:\t%s\n", service.HealthCheck.Protocol)
				fmt.Fprintf(out, "      Timeout:\t%vs\n", service.HealthCheck.Timeout.Seconds())
				fmt.Fprintf(out, "      Interval:\tevery %vs\n", service.HealthCheck.Interval.Seconds())
				fmt.Fprintf(out, "      Retries:\t%d\n", service.HealthCheck.Retries)
				if service.HealthCheck.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					fmt.Fprintf(out, "      HTTP Domain:\t%s\n", service.HealthCheck.HTTP.Domain)
					fmt.Fprintf(out, "      HTTP Path:\t%s\n", service.HealthCheck.HTTP.Path)
					fmt.Fprintf(out, "      Response:\t%s\n", service.HealthCheck.HTTP.Response)
					fmt.Fprintf(out, "      TLS:\t%s\n", util.YesNo(service.HealthCheck.HTTP.TLS))
					fmt.Fprintf(out, "      Status Codes:\t%v\n", service.HealthCheck.HTTP.StatusCodes)
				}
			}
		}

		if len(loadBalancer.Targets) == 0 {
			fmt.Fprintf(out, "Targets:\tNo targets\n")
		} else {
			fmt.Fprintf(out, "Targets:\t\n")
			for _, target := range loadBalancer.Targets {
				fmt.Fprintf(out, "  - Type:\t%s\n", target.Type)
				switch target.Type {
				case hcloud.LoadBalancerTargetTypeServer:
					fmt.Fprintf(out, "    Server:\t\n")
					fmt.Fprintf(out, "      ID:\t%d\n", target.Server.Server.ID)
					fmt.Fprintf(out, "      Name:\t%s\n", s.Client().Server().ServerName(target.Server.Server.ID))
					fmt.Fprintf(out, "    Use Private IP:\t%s\n", util.YesNo(target.UsePrivateIP))
					fmt.Fprintf(out, "    Status:\t\n")
					for _, healthStatus := range target.HealthStatus {
						fmt.Fprintf(out, "    - Service:\t%d\n", healthStatus.ListenPort)
						fmt.Fprintf(out, "      Status:\t%s\n", healthStatus.Status)
					}
				case hcloud.LoadBalancerTargetTypeLabelSelector:
					fmt.Fprintf(out, "    Label Selector:\t%s\n", target.LabelSelector.Selector)
					fmt.Fprintf(out, "      Targets: (%d)\t\n", len(target.Targets))
					if len(target.Targets) == 0 {
						fmt.Fprintf(out, "      No targets\t\n")
					}
					if !withLabelSelectorTargets {
						continue
					}
					for _, lbtarget := range target.Targets {
						fmt.Fprintf(out, "      - Type:\t%s\n", lbtarget.Type)
						fmt.Fprintf(out, "        Server ID:\t%d\n", lbtarget.Server.Server.ID)
						fmt.Fprintf(out, "        Status:\t\n")
						for _, healthStatus := range lbtarget.HealthStatus {
							fmt.Fprintf(out, "          - Service:\t%d\n", healthStatus.ListenPort)
							fmt.Fprintf(out, "            Status:\t%s\n", healthStatus.Status)
						}
					}
				case hcloud.LoadBalancerTargetTypeIP:
					fmt.Fprintf(out, "    IP:\t%s\n", target.IP.IP)
					fmt.Fprintf(out, "    Status:\t\n")
					for _, healthStatus := range target.HealthStatus {
						fmt.Fprintf(out, "    - Service:\t%d\n", healthStatus.ListenPort)
						fmt.Fprintf(out, "      Status:\t%s\n", healthStatus.Status)
					}
				}
			}
		}

		fmt.Fprintf(out, "Traffic:\t\n")
		fmt.Fprintf(out, "  Outgoing:\t%v\n", humanize.IBytes(loadBalancer.OutgoingTraffic))
		fmt.Fprintf(out, "  Ingoing:\t%v\n", humanize.IBytes(loadBalancer.IngoingTraffic))
		fmt.Fprintf(out, "  Included:\t%v\n", humanize.IBytes(loadBalancer.IncludedTraffic))

		fmt.Fprintf(out, "Protection:\t\n")
		fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(loadBalancer.Protection.Delete))

		util.DescribeLabels(out, loadBalancer.Labels, "")

		return nil
	},
}
