package firewall

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var AddRuleCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "add-rule [options] (--direction in --source-ips <ips> | --direction out --destination-ips <ips>) (--protocol <tcp|udp> --port <port> | --protocol <icmp|esp|gre>) <firewall>",
			Short:                 "Add a single rule to a firewall",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Firewall().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("direction", "", "Direction (in, out) (required)")
		_ = cmd.RegisterFlagCompletionFunc("direction", cmpl.SuggestCandidates("in", "out"))
		_ = cmd.MarkFlagRequired("direction")

		cmd.Flags().String("protocol", "", "Protocol (icmp, esp, gre, udp or tcp) (required)")
		_ = cmd.RegisterFlagCompletionFunc("protocol", cmpl.SuggestCandidates("icmp", "udp", "tcp", "esp", "gre"))
		_ = cmd.MarkFlagRequired("protocol")

		cmd.Flags().StringSlice("source-ips", []string{}, "Source IPs (CIDR Notation) (required when direction is in)")

		cmd.Flags().StringSlice("destination-ips", []string{}, "Destination IPs (CIDR Notation) (required when direction is out)")

		cmd.Flags().String("port", "", "Port to which traffic will be allowed, only applicable for protocols TCP and UDP, you can specify port ranges, sample: 80-85")

		cmd.Flags().String("description", "", "Description of the Firewall rule")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		firewall, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return err
		}
		if firewall == nil {
			return fmt.Errorf("Firewall not found: %v", idOrName)
		}

		rule, err := parseRuleFromArgs(cmd.Flags())
		if err != nil {
			return err
		}

		rules := append(firewall.Rules, *rule)

		actions, _, err := s.Client().Firewall().SetRules(s, firewall,
			hcloud.FirewallSetRulesOpts{Rules: rules},
		)
		if err != nil {
			return err
		}
		if err := s.WaitForActions(s, cmd, actions...); err != nil {
			return err
		}

		cmd.Printf("Firewall Rules for Firewall %d updated\n", firewall.ID)

		return nil
	},
}

func parseRuleFromArgs(flags *pflag.FlagSet) (*hcloud.FirewallRule, error) {
	direction, _ := flags.GetString("direction")
	protocol, _ := flags.GetString("protocol")
	sourceIPs, _ := flags.GetStringSlice("source-ips")
	destinationIPs, _ := flags.GetStringSlice("destination-ips")
	port, _ := flags.GetString("port")
	description, _ := flags.GetString("description")

	rule := &hcloud.FirewallRule{
		SourceIPs:      make([]net.IPNet, 0),
		DestinationIPs: make([]net.IPNet, 0),
	}

	switch hcloud.FirewallRuleDirection(direction) {
	case hcloud.FirewallRuleDirectionIn, hcloud.FirewallRuleDirectionOut:
		rule.Direction = hcloud.FirewallRuleDirection(direction)
	default:
		return nil, fmt.Errorf("invalid direction: %s", direction)
	}

	switch hcloud.FirewallRuleProtocol(protocol) {
	case hcloud.FirewallRuleProtocolTCP, hcloud.FirewallRuleProtocolUDP, hcloud.FirewallRuleProtocolICMP, hcloud.FirewallRuleProtocolESP, hcloud.FirewallRuleProtocolGRE:
		rule.Protocol = hcloud.FirewallRuleProtocol(protocol)
	default:
		return nil, fmt.Errorf("invalid protocol: %s", protocol)
	}

	if port != "" {
		rule.Port = hcloud.Ptr(port)
	}

	if description != "" {
		rule.Description = hcloud.Ptr(description)
	}

	switch rule.Protocol {
	case hcloud.FirewallRuleProtocolUDP, hcloud.FirewallRuleProtocolTCP:
		if port == "" {
			return nil, fmt.Errorf("port is required (--port)")
		}
	default:
		if port != "" {
			return nil, fmt.Errorf("port is not allowed for this protocol")
		}
	}

	switch rule.Direction {
	case hcloud.FirewallRuleDirectionOut:
		for i, ip := range destinationIPs {
			n, err := ValidateFirewallIP(ip)
			if err != nil {
				return nil, fmt.Errorf("destination error on index %d: %w", i, err)
			}
			rule.DestinationIPs = append(rule.DestinationIPs, *n)
		}
	case hcloud.FirewallRuleDirectionIn:
		for i, ip := range sourceIPs {
			n, err := ValidateFirewallIP(ip)
			if err != nil {
				return nil, fmt.Errorf("source ips error on index %d: %w", i, err)
			}
			rule.SourceIPs = append(rule.SourceIPs, *n)
		}
	}

	return rule, nil
}
