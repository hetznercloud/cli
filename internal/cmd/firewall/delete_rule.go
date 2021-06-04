package firewall

import (
	"fmt"
	"net"
	"reflect"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func newDeleteRuleCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete-rule FIREWALL FLAGS",
		Short:                 "Delete a single rule to a firewall",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FirewallNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(cli.EnsureToken),
		RunE:                  cli.Wrap(runDeleteRule),
	}
	cmd.Flags().String("direction", "", "Direction (in, out) (required)")
	cmd.RegisterFlagCompletionFunc("direction", cmpl.SuggestCandidates("in", "out"))
	cmd.MarkFlagRequired("direction")

	cmd.Flags().String("protocol", "", "Protocol (icmp, udp or tcp) (required)")
	cmd.RegisterFlagCompletionFunc("protocol", cmpl.SuggestCandidates("icmp", "udp", "tcp"))
	cmd.MarkFlagRequired("protocol")

	cmd.Flags().StringArray("source-ips", []string{}, "Source IPs (CIDR Notation) (required when direction is in)")

	cmd.Flags().StringArray("destination-ips", []string{}, "Destination IPs (CIDR Notation) (required when direction is out)")

	cmd.Flags().String("port", "", "Port to which traffic will be allowed, only applicable for protocols TCP and UDP")
	return cmd
}

func runDeleteRule(cli *state.State, cmd *cobra.Command, args []string) error {
	direction, _ := cmd.Flags().GetString("direction")
	protocol, _ := cmd.Flags().GetString("protocol")
	sourceIPs, _ := cmd.Flags().GetStringArray("source-ips")
	destinationIPs, _ := cmd.Flags().GetStringArray("destination-ips")
	port, _ := cmd.Flags().GetString("port")

	idOrName := args[0]
	firewall, _, err := cli.Client().Firewall.Get(cli.Context, idOrName)
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

	switch rule.Protocol {
	case hcloud.FirewallRuleProtocolTCP:
	case hcloud.FirewallRuleProtocolUDP:
		if port == "" {
			return fmt.Errorf("port is required")
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
			n, err := validateFirewallIP(ip)
			if err != nil {
				return fmt.Errorf("destination ips error on index %d: %s", i, err)
			}
			rule.DestinationIPs[i] = *n
			rule.SourceIPs = make([]net.IPNet, 0)
		}
	case hcloud.FirewallRuleDirectionIn:
		rule.SourceIPs = make([]net.IPNet, len(sourceIPs))
		for i, ip := range sourceIPs {
			n, err := validateFirewallIP(ip)
			if err != nil {
				return fmt.Errorf("source ips error on index %d: %s", i, err)
			}
			rule.DestinationIPs = make([]net.IPNet, 0)
			rule.SourceIPs[i] = *n
		}
	}

	var rules []hcloud.FirewallRule
	for _, existingRule := range firewall.Rules {
		if !reflect.DeepEqual(existingRule, rule) {
			rules = append(rules, existingRule)
		}
	}
	if len(rules) == len(firewall.Rules) {
		return fmt.Errorf("the specified rule was not found in the ruleset of Firewall %d", firewall.ID)
	}
	actions, _, err := cli.Client().Firewall.SetRules(cli.Context, firewall,
		hcloud.FirewallSetRulesOpts{Rules: rules},
	)
	if err != nil {
		return err
	}
	if err := cli.ActionsProgresses(cli.Context, actions); err != nil {
		return err
	}
	fmt.Printf("Firewall Rules for Firewall %d updated\n", firewall.ID)

	return nil
}
