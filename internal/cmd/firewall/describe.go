package firewall

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "firewall",
	ShortDescription:     "Describe an firewall",
	JSONKeyGetByID:       "firewall",
	JSONKeyGetByName:     "firewalls",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Firewall().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		firewall := resource.(*hcloud.Firewall)

		fmt.Printf("ID:\t\t%d\n", firewall.ID)
		fmt.Printf("Name:\t\t%s\n", firewall.Name)
		fmt.Printf("Created:\t%s (%s)\n", util.Datetime(firewall.Created), humanize.Time(firewall.Created))

		fmt.Print("Labels:\n")
		if len(firewall.Labels) == 0 {
			fmt.Print("  No labels\n")
		} else {
			for key, value := range firewall.Labels {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		fmt.Print("Rules:\n")
		if len(firewall.Rules) == 0 {
			fmt.Print("  No rules\n")
		} else {
			for _, rule := range firewall.Rules {
				fmt.Printf("  - Direction:\t\t%s\n", rule.Direction)
				if rule.Description != nil {
					fmt.Printf("    Description:\t%s\n", *rule.Description)
				}
				fmt.Printf("    Protocol:\t\t%s\n", rule.Protocol)
				if rule.Port != nil {
					fmt.Printf("    Port:\t\t%s\n", *rule.Port)
				}

				var ips []net.IPNet
				switch rule.Direction {
				case hcloud.FirewallRuleDirectionIn:
					fmt.Print("    Source IPs:\n")
					ips = rule.SourceIPs
				case hcloud.FirewallRuleDirectionOut:
					fmt.Print("    Destination IPs:\n")
					ips = rule.DestinationIPs
				}

				for _, cidr := range ips {
					fmt.Printf("     \t\t\t%s\n", cidr.String())
				}
			}
		}
		fmt.Print("Applied To:\n")
		if len(firewall.AppliedTo) == 0 {
			fmt.Print("  Not applied\n")
		} else {
			for _, resource := range firewall.AppliedTo {
				fmt.Printf("  - Type:\t\t%s\n", resource.Type)
				switch resource.Type {
				case hcloud.FirewallResourceTypeServer:
					fmt.Printf("    Server ID:\t\t%d\n", resource.Server.ID)
					fmt.Printf("    Server Name:\t%s\n", client.Server().ServerName(resource.Server.ID))
				case hcloud.FirewallResourceTypeLabelSelector:
					fmt.Printf("    Label Selector:\t%s\n", resource.LabelSelector.Selector)
				}
			}
		}
		return nil
	},
}
