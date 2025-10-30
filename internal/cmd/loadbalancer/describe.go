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

		_, _ = fmt.Fprintf(out, "ID:\t%d\n", loadBalancer.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", loadBalancer.Name)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(loadBalancer.Created), humanize.Time(loadBalancer.Created))
		_, _ = fmt.Fprintf(out, "Public Net:\t\n")
		_, _ = fmt.Fprintf(out, "  Enabled:\t%s\n", util.YesNo(loadBalancer.PublicNet.Enabled))
		_, _ = fmt.Fprintf(out, "  IPv4:\t%s\n", loadBalancer.PublicNet.IPv4.IP.String())
		_, _ = fmt.Fprintf(out, "  IPv4 DNS PTR:\t%s\n", loadBalancer.PublicNet.IPv4.DNSPtr)
		_, _ = fmt.Fprintf(out, "  IPv6:\t%s\n", loadBalancer.PublicNet.IPv6.IP.String())
		_, _ = fmt.Fprintf(out, "  IPv6 DNS PTR:\t%s\n", loadBalancer.PublicNet.IPv6.DNSPtr)

		if len(loadBalancer.PrivateNet) > 0 {
			_, _ = fmt.Fprintf(out, "Private Net:\t\n")
			for _, n := range loadBalancer.PrivateNet {
				_, _ = fmt.Fprintf(out, "  - ID:\t%d\n", n.Network.ID)
				_, _ = fmt.Fprintf(out, "    Name:\t%s\n", s.Client().Network().Name(n.Network.ID))
				_, _ = fmt.Fprintf(out, "    IP:\t%s\n", n.IP.String())
			}
		} else {
			_, _ = fmt.Fprintf(out, "Private Net:\tNo Private Network\n")
		}
		_, _ = fmt.Fprintf(out, "Algorithm:\t%s\n", loadBalancer.Algorithm.Type)

		_, _ = fmt.Fprintf(out, "Load Balancer Type:\t\n")
		_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(loadbalancertype.DescribeLoadBalancerType(s, loadBalancer.LoadBalancerType, true), "  "))

		if len(loadBalancer.Services) == 0 {
			_, _ = fmt.Fprintf(out, "Services:\tNo services\n")
		} else {
			_, _ = fmt.Fprintf(out, "Services:\t\n")
			for _, service := range loadBalancer.Services {
				_, _ = fmt.Fprintf(out, "  - Protocol:\t%s\n", service.Protocol)
				_, _ = fmt.Fprintf(out, "    Listen Port:\t%d\n", service.ListenPort)
				_, _ = fmt.Fprintf(out, "    Destination Port:\t%d\n", service.DestinationPort)
				_, _ = fmt.Fprintf(out, "    Proxy Protocol:\t%s\n", util.YesNo(service.Proxyprotocol))
				if service.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					_, _ = fmt.Fprintf(out, "    Sticky Sessions:\t%s\n", util.YesNo(service.HTTP.StickySessions))
					if service.HTTP.StickySessions {
						_, _ = fmt.Fprintf(out, "    Sticky Cookie Name:\t%s\n", service.HTTP.CookieName)
						_, _ = fmt.Fprintf(out, "    Sticky Cookie Lifetime:\t%vs\n", service.HTTP.CookieLifetime.Seconds())
					}
					if service.Protocol == hcloud.LoadBalancerServiceProtocolHTTPS {
						_, _ = fmt.Fprintf(out, "    Certificates:\n")
						for _, cert := range service.HTTP.Certificates {
							_, _ = fmt.Fprintf(out, "      - ID:\t%v\n", cert.ID)
						}
					}
				}

				_, _ = fmt.Fprintf(out, "    Health Check:\n")
				_, _ = fmt.Fprintf(out, "      Protocol:\t%s\n", service.HealthCheck.Protocol)
				_, _ = fmt.Fprintf(out, "      Timeout:\t%vs\n", service.HealthCheck.Timeout.Seconds())
				_, _ = fmt.Fprintf(out, "      Interval:\tevery %vs\n", service.HealthCheck.Interval.Seconds())
				_, _ = fmt.Fprintf(out, "      Retries:\t%d\n", service.HealthCheck.Retries)
				if service.HealthCheck.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
					_, _ = fmt.Fprintf(out, "      HTTP Domain:\t%s\n", service.HealthCheck.HTTP.Domain)
					_, _ = fmt.Fprintf(out, "      HTTP Path:\t%s\n", service.HealthCheck.HTTP.Path)
					_, _ = fmt.Fprintf(out, "      Response:\t%s\n", service.HealthCheck.HTTP.Response)
					_, _ = fmt.Fprintf(out, "      TLS:\t%s\n", util.YesNo(service.HealthCheck.HTTP.TLS))
					_, _ = fmt.Fprintf(out, "      Status Codes:\t%v\n", service.HealthCheck.HTTP.StatusCodes)
				}
			}
		}

		if len(loadBalancer.Targets) == 0 {
			_, _ = fmt.Fprintf(out, "Targets:\tNo targets\n")
		} else {
			_, _ = fmt.Fprintf(out, "Targets:\t\n")
			for _, target := range loadBalancer.Targets {
				_, _ = fmt.Fprintf(out, "  - Type:\t%s\n", target.Type)
				switch target.Type {
				case hcloud.LoadBalancerTargetTypeServer:
					_, _ = fmt.Fprintf(out, "    Server:\t\n")
					_, _ = fmt.Fprintf(out, "      ID:\t%d\n", target.Server.Server.ID)
					_, _ = fmt.Fprintf(out, "      Name:\t%s\n", s.Client().Server().ServerName(target.Server.Server.ID))
					_, _ = fmt.Fprintf(out, "    Use Private IP:\t%s\n", util.YesNo(target.UsePrivateIP))
					_, _ = fmt.Fprintf(out, "    Status:\t\n")
					for _, healthStatus := range target.HealthStatus {
						_, _ = fmt.Fprintf(out, "    - Service:\t%d\n", healthStatus.ListenPort)
						_, _ = fmt.Fprintf(out, "      Status:\t%s\n", healthStatus.Status)
					}
				case hcloud.LoadBalancerTargetTypeLabelSelector:
					_, _ = fmt.Fprintf(out, "    Label Selector:\t%s\n", target.LabelSelector.Selector)
					_, _ = fmt.Fprintf(out, "      Targets: (%d)\t\n", len(target.Targets))
					if len(target.Targets) == 0 {
						_, _ = fmt.Fprintf(out, "      No targets\t\n")
					}
					if !withLabelSelectorTargets {
						continue
					}
					for _, lbtarget := range target.Targets {
						_, _ = fmt.Fprintf(out, "      - Type:\t%s\n", lbtarget.Type)
						_, _ = fmt.Fprintf(out, "        Server ID:\t%d\n", lbtarget.Server.Server.ID)
						_, _ = fmt.Fprintf(out, "        Status:\t\n")
						for _, healthStatus := range lbtarget.HealthStatus {
							_, _ = fmt.Fprintf(out, "          - Service:\t%d\n", healthStatus.ListenPort)
							_, _ = fmt.Fprintf(out, "            Status:\t%s\n", healthStatus.Status)
						}
					}
				case hcloud.LoadBalancerTargetTypeIP:
					_, _ = fmt.Fprintf(out, "    IP:\t%s\n", target.IP.IP)
					_, _ = fmt.Fprintf(out, "    Status:\t\n")
					for _, healthStatus := range target.HealthStatus {
						_, _ = fmt.Fprintf(out, "    - Service:\t%d\n", healthStatus.ListenPort)
						_, _ = fmt.Fprintf(out, "      Status:\t%s\n", healthStatus.Status)
					}
				}
			}
		}

		_, _ = fmt.Fprintf(out, "Traffic:\t\n")
		_, _ = fmt.Fprintf(out, "  Outgoing:\t%v\n", humanize.IBytes(loadBalancer.OutgoingTraffic))
		_, _ = fmt.Fprintf(out, "  Ingoing:\t%v\n", humanize.IBytes(loadBalancer.IngoingTraffic))
		_, _ = fmt.Fprintf(out, "  Included:\t%v\n", humanize.IBytes(loadBalancer.IncludedTraffic))

		_, _ = fmt.Fprintf(out, "Protection:\t\n")
		_, _ = fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(loadBalancer.Protection.Delete))

		util.DescribeLabels(out, loadBalancer.Labels, "")

		return nil
	},
}
