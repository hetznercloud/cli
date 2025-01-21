package firewall

import (
	"fmt"

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
		_ = cmd.RegisterFlagCompletionFunc("direction", cmpl.SuggestCandidates("in", "out"))
		_ = cmd.MarkFlagRequired("direction")

		cmd.Flags().String("protocol", "", "Protocol (icmp, esp, gre, udp or tcp) (required)")
		_ = cmd.RegisterFlagCompletionFunc("protocol", cmpl.SuggestCandidates("icmp", "udp", "tcp", "esp", "gre"))
		_ = cmd.MarkFlagRequired("protocol")

		cmd.Flags().StringArray("source-ips", []string{}, "Source IPs (CIDR Notation) (required when direction is in)")

		cmd.Flags().StringArray("destination-ips", []string{}, "Destination IPs (CIDR Notation) (required when direction is out)")

		cmd.Flags().String("port", "", "Port to which traffic will be allowed, only applicable for protocols TCP and UDP")

		cmd.Flags().String("description", "", "Description of the firewall rule")
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

		var rules = make([]hcloud.FirewallRule, 0)
		for _, existingRule := range firewall.Rules {
			if !EqualFirewallRule(existingRule, *rule) {
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
		if err := s.WaitForActions(s, cmd, actions...); err != nil {
			return err
		}
		cmd.Printf("Firewall Rules for Firewall %d updated\n", firewall.ID)

		return nil
	},
}
