package firewall

import (
	"fmt"
	"net"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteRuleCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "delete-rule [options] (--direction in --source-ips <ips> | --direction out --destination-ips <ips>) (--protocol <tcp|udp> --port <port> | --protocol <icmp|esp|gre>) <firewall>",
			Short:                 "Delete a single rule to a firewall",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Firewall().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("direction", "", "Direction (in, out) (required)")
		cmd.RegisterFlagCompletionFunc("direction", cmpl.SuggestCandidates("in", "out"))
		cmd.MarkFlagRequired("direction")

		cmd.Flags().String("protocol", "", "Protocol (icmp, esp, gre, udp or tcp) (required)")
		cmd.RegisterFlagCompletionFunc("protocol", cmpl.SuggestCandidates("icmp", "udp", "tcp", "esp", "gre"))
		cmd.MarkFlagRequired("protocol")

		cmd.Flags().StringArray("source-ips", []string{}, "Source IPs (CIDR Notation) (required when direction is in)")

		cmd.Flags().StringArray("destination-ips", []string{}, "Destination IPs (CIDR Notation) (required when direction is out)")

		cmd.Flags().String("port", "", "Port to which traffic will be allowed, only applicable for protocols TCP and UDP")

		cmd.Flags().String("description", "", "Description of the firewall rule")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		direction, _ := cmd.Flags().GetString("direction")
		protocol, _ := cmd.Flags().GetString("protocol")
		sourceIPs, _ := cmd.Flags().GetStringArray("source-ips")
		destinationIPs, _ := cmd.Flags().GetStringArray("destination-ips")
		port, _ := cmd.Flags().GetString("port")
		description, _ := cmd.Flags().GetString("description")

		idOrName := args[0]
		firewall, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return err
		}
		if firewall == nil {
			return fmt.Errorf("Firewall not found: %v", idOrName)
		}

		d := hcloud.FirewallRuleDirection(direction)
		rule := hcloud.FirewallRule{
			Direction: d,
			Protocol:  hcloud.FirewallRuleProtocol(protocol),
		}
		if port != "" {
			rule.Port = hcloud.String(port)
		}
		if description != "" {
			rule.Description = hcloud.String(description)
		}

		switch rule.Protocol {
		case hcloud.FirewallRuleProtocolTCP, hcloud.FirewallRuleProtocolUDP:
			if port == "" {
				return fmt.Errorf("port is required (--port)")
			}
		default:
			if port != "" {
				return fmt.Errorf("port is not allowed for this protocol")
			}
		}

		switch d {
		case hcloud.FirewallRuleDirectionOut:
			rule.DestinationIPs = make([]net.IPNet, len(destinationIPs))
			for i, ip := range destinationIPs {
				n, err := ValidateFirewallIP(ip)
				if err != nil {
					return fmt.Errorf("destination ips error on index %d: %s", i, err)
				}
				rule.DestinationIPs[i] = *n
				rule.SourceIPs = make([]net.IPNet, 0)
			}
		case hcloud.FirewallRuleDirectionIn:
			rule.SourceIPs = make([]net.IPNet, len(sourceIPs))
			for i, ip := range sourceIPs {
				n, err := ValidateFirewallIP(ip)
				if err != nil {
					return fmt.Errorf("source ips error on index %d: %s", i, err)
				}
				rule.DestinationIPs = make([]net.IPNet, 0)
				rule.SourceIPs[i] = *n
			}
		}

		var rules = make([]hcloud.FirewallRule, 0)
		for _, existingRule := range firewall.Rules {
			if !reflect.DeepEqual(existingRule, rule) {
				rules = append(rules, existingRule)
			}
		}
		if len(rules) == len(firewall.Rules) {
			return fmt.Errorf("the specified rule was not found in the ruleset of Firewall %d", firewall.ID)
		}
		actions, _, err := s.Client().Firewall().SetRules(s, firewall,
			hcloud.FirewallSetRulesOpts{Rules: rules},
		)
		if err != nil {
			return err
		}
		if err := s.WaitForActions(cmd, s, actions...); err != nil {
			return err
		}
		cmd.Printf("Firewall Rules for Firewall %d updated\n", firewall.ID)

		return nil
	},
}
