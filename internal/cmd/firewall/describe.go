package firewall

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Firewall]{
	ResourceNameSingular: "Firewall",
	ShortDescription:     "Describe a Firewall",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Firewall, any, error) {
		fw, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return fw, hcloud.SchemaFromFirewall(fw), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, firewall *hcloud.Firewall) error {
		_, _ = fmt.Fprintf(out, "ID:\t%d\n", firewall.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", firewall.Name)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(firewall.Created), humanize.Time(firewall.Created))

		util.DescribeLabels(out, firewall.Labels, "")

		if len(firewall.Rules) == 0 {
			_, _ = fmt.Fprintf(out, "Rules:\tNo rules\n")
		} else {
			_, _ = fmt.Fprintf(out, "Rules:\t\n")
			for _, rule := range firewall.Rules {
				_, _ = fmt.Fprintf(out, "  - Direction:\t%s\n", rule.Direction)
				if rule.Description != nil {
					_, _ = fmt.Fprintf(out, "    Description:\t%s\n", *rule.Description)
				}
				_, _ = fmt.Fprintf(out, "    Protocol:\t%s\n", rule.Protocol)
				if rule.Port != nil {
					_, _ = fmt.Fprintf(out, "    Port:\t%s\n", *rule.Port)
				}

				var ips []net.IPNet
				switch rule.Direction {
				case hcloud.FirewallRuleDirectionIn:
					_, _ = fmt.Fprintf(out, "    Source IPs:\t\n")
					ips = rule.SourceIPs
				case hcloud.FirewallRuleDirectionOut:
					_, _ = fmt.Fprintf(out, "    Destination IPs:\t\n")
					ips = rule.DestinationIPs
				}

				for _, cidr := range ips {
					_, _ = fmt.Fprintf(out, "\t%s\n", cidr.String())
				}
			}
		}

		_, _ = fmt.Fprintf(out, "\n")

		if len(firewall.AppliedTo) == 0 {
			_, _ = fmt.Fprintf(out, "Applied To:\nNot applied\n")
		} else {
			_, _ = fmt.Fprintf(out, "Applied To:\t\n")
			_, _ = fmt.Fprintf(out, "%s", describeResources(s.Client(), firewall.AppliedTo))
		}

		return nil
	},
}

func describeResources(client hcapi2.Client, resources []hcloud.FirewallResource) string {
	var sb strings.Builder

	for _, resource := range resources {
		sb.WriteString(fmt.Sprintf("  - Type:\t%s\n", resource.Type))

		switch resource.Type {
		case hcloud.FirewallResourceTypeServer:
			sb.WriteString(fmt.Sprintf("    Server ID:\t%d\n", resource.Server.ID))
			sb.WriteString(fmt.Sprintf("    Server Name:\t%s\n", client.Server().ServerName(resource.Server.ID)))

		case hcloud.FirewallResourceTypeLabelSelector:
			sb.WriteString(fmt.Sprintf("    Label Selector:\t%s\n", resource.LabelSelector.Selector))
			if len(resource.AppliedToResources) > 0 {
				sb.WriteString("    Applied to resources:\t\n")
				substr := describeResources(client, resource.AppliedToResources)
				sb.WriteString(util.PrefixLines(substr, "  "))
			}
		}
	}

	return sb.String()
}
