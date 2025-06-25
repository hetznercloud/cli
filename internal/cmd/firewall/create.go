package firewall

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create [options] --name <name>",
			Short: "Create a Firewall",
		}
		cmd.Flags().String("name", "", "Name")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().String("rules-file", "", "JSON file containing your routes (use - to read from stdin). The structure of the file needs to be the same as within the API: https://docs.hetzner.cloud/reference/cloud#firewalls-get-a-firewall ")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		labels, _ := cmd.Flags().GetStringToString("label")

		opts := hcloud.FirewallCreateOpts{
			Name:   name,
			Labels: labels,
		}

		rulesFile, _ := cmd.Flags().GetString("rules-file")
		if rulesFile != "" {
			rules, err := parseRulesFile(rulesFile)
			if err != nil {
				return nil, nil, err
			}
			opts.Rules = rules
		}

		result, _, err := s.Client().Firewall().Create(s, opts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(s, cmd, result.Actions...); err != nil {
			return nil, nil, err
		}

		cmd.Printf("Firewall %d created\n", result.Firewall.ID)

		return result.Firewall, util.Wrap("firewall", hcloud.SchemaFromFirewall(result.Firewall)), err
	},
}

func parseRulesFile(path string) ([]hcloud.FirewallRule, error) {
	var (
		data []byte
		err  error
	)
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}

	var ruleSchemas []schema.FirewallRule
	err = json.Unmarshal(data, &ruleSchemas)
	if err != nil {
		return nil, err
	}

	rules := make([]hcloud.FirewallRule, 0, len(ruleSchemas))
	for _, rule := range ruleSchemas {
		var sourceNets []net.IPNet
		for i, sourceIP := range rule.SourceIPs {
			_, sourceNet, err := net.ParseCIDR(sourceIP)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR on index %d : %w", i, err)
			}
			sourceNets = append(sourceNets, *sourceNet)
		}
		var destNets []net.IPNet
		for i, destIP := range rule.DestinationIPs {
			_, destNet, err := net.ParseCIDR(destIP)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR on index %d : %w", i, err)
			}
			destNets = append(destNets, *destNet)
		}
		rules = append(rules, hcloud.FirewallRule{
			Direction:      hcloud.FirewallRuleDirection(rule.Direction),
			SourceIPs:      sourceNets,
			DestinationIPs: destNets,
			Protocol:       hcloud.FirewallRuleProtocol(rule.Protocol),
			Port:           rule.Port,
			Description:    rule.Description,
		})
	}
	return rules, nil
}
