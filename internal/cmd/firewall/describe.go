package firewall

import (
	"context"
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
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Firewall().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Firewall, any, error) {
		fw, _, err := s.Client().Firewall().Get(cmd.Context(), idOrName)
		if err != nil {
			return nil, nil, err
		}
		return fw, hcloud.SchemaFromFirewall(fw), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, out io.Writer, firewall *hcloud.Firewall) error {
		fmt.Fprintf(out, "ID:\t%d\n", firewall.ID)
		fmt.Fprintf(out, "Name:\t%s\n", firewall.Name)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(firewall.Created), humanize.Time(firewall.Created))

		fmt.Fprintln(out)
		util.DescribeLabels(out, firewall.Labels, "")

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Rules:\n")
		if len(firewall.Rules) == 0 {
			fmt.Fprintf(out, "  No rules\n")
		} else {
			for i, rule := range firewall.Rules {
				if i > 0 {
					fmt.Fprintln(out)
				}

				fmt.Fprintf(out, "  - Direction:\t%s\n", rule.Direction)
				if rule.Description != nil {
					fmt.Fprintf(out, "    Description:\t%s\n", *rule.Description)
				}
				fmt.Fprintf(out, "    Protocol:\t%s\n", rule.Protocol)
				if rule.Port != nil {
					fmt.Fprintf(out, "    Port:\t%s\n", *rule.Port)
				}

				var ips []net.IPNet
				switch rule.Direction {
				case hcloud.FirewallRuleDirectionIn:
					fmt.Fprintf(out, "    Source IPs:\n")
					ips = rule.SourceIPs
				case hcloud.FirewallRuleDirectionOut:
					fmt.Fprintf(out, "    Destination IPs:\n")
					ips = rule.DestinationIPs
				}

				for _, cidr := range ips {
					fmt.Fprintf(out, "      %s\n", cidr.String())
				}
			}
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Applied To:\n")
		if len(firewall.AppliedTo) == 0 {
			fmt.Fprintf(out, "  Not applied\n")
		} else {
			resources, err := describeResources(cmd.Context(), s.Client(), firewall.AppliedTo)
			if err != nil {
				return err
			}
			fmt.Fprint(out, resources)
		}

		return nil
	},
}

func describeResources(ctx context.Context, client hcapi2.Client, resources []hcloud.FirewallResource) (string, error) {
	var sb strings.Builder

	for _, resource := range resources {
		fmt.Fprintf(&sb, "  - Type:\t%s\n", resource.Type)

		switch resource.Type {
		case hcloud.FirewallResourceTypeServer:
			name, err := client.Server().ServerName(ctx, resource.Server.ID)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&sb, "    Server ID:\t%d\n", resource.Server.ID)
			fmt.Fprintf(&sb, "    Server Name:\t%s\n", name)

		case hcloud.FirewallResourceTypeLabelSelector:
			fmt.Fprintf(&sb, "    Label Selector:\t%s\n", resource.LabelSelector.Selector)

			if len(resource.AppliedToResources) > 0 {
				fmt.Fprintf(&sb, "    Applied to resources:\n")
				substr, err := describeResources(ctx, client, resource.AppliedToResources)
				if err != nil {
					return "", err
				}
				fmt.Fprint(&sb, util.PrefixLines(substr, "  "))
			}
		}
	}

	return sb.String(), nil
}
