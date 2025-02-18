package firewall

import (
	"net"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Firewall]{
	ResourceNameSingular: "firewall",
	ShortDescription:     "Describe an firewall",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Firewall, any, error) {
		fw, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return fw, hcloud.SchemaFromFirewall(fw), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, firewall *hcloud.Firewall) error {
		cmd.Printf("ID:\t\t%d\n", firewall.ID)
		cmd.Printf("Name:\t\t%s\n", firewall.Name)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(firewall.Created), humanize.Time(firewall.Created))

		cmd.Print("Labels:\n")
		if len(firewall.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range util.IterateInOrder(firewall.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		cmd.Print("Rules:\n")
		if len(firewall.Rules) == 0 {
			cmd.Print("  No rules\n")
		} else {
			for _, rule := range firewall.Rules {
				cmd.Printf("  - Direction:\t\t%s\n", rule.Direction)
				if rule.Description != nil {
					cmd.Printf("    Description:\t%s\n", *rule.Description)
				}
				cmd.Printf("    Protocol:\t\t%s\n", rule.Protocol)
				if rule.Port != nil {
					cmd.Printf("    Port:\t\t%s\n", *rule.Port)
				}

				var ips []net.IPNet
				switch rule.Direction {
				case hcloud.FirewallRuleDirectionIn:
					cmd.Print("    Source IPs:\n")
					ips = rule.SourceIPs
				case hcloud.FirewallRuleDirectionOut:
					cmd.Print("    Destination IPs:\n")
					ips = rule.DestinationIPs
				}

				for _, cidr := range ips {
					cmd.Printf("     \t\t\t%s\n", cidr.String())
				}
			}
		}
		cmd.Print("Applied To:\n")
		if len(firewall.AppliedTo) == 0 {
			cmd.Print("  Not applied\n")
		} else {
			for _, resource := range firewall.AppliedTo {
				cmd.Printf("  - Type:\t\t%s\n", resource.Type)
				switch resource.Type {
				case hcloud.FirewallResourceTypeServer:
					cmd.Printf("    Server ID:\t\t%d\n", resource.Server.ID)
					cmd.Printf("    Server Name:\t%s\n", s.Client().Server().ServerName(resource.Server.ID))
				case hcloud.FirewallResourceTypeLabelSelector:
					cmd.Printf("    Label Selector:\t%s\n", resource.LabelSelector.Selector)
				}
			}
		}
		return nil
	},
}
