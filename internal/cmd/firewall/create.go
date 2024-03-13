package firewall

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create [options] --name <name>",
			Short: "Create a Firewall",
		}
		cmd.Flags().String("name", "", "Name")
		cmd.MarkFlagRequired("name")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().String("rules-file", "", "JSON file containing your routes (use - to read from stdin). The structure of the file needs to be the same as within the API: https://docs.hetzner.cloud/#firewalls-get-a-firewall ")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, strings []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		labels, _ := cmd.Flags().GetStringToString("label")

		opts := hcloud.FirewallCreateOpts{
			Name:   name,
			Labels: labels,
		}

		rulesFile, _ := cmd.Flags().GetString("rules-file")

		if len(rulesFile) > 0 {
			var data []byte
			var err error
			if rulesFile == "-" {
				data, err = ioutil.ReadAll(os.Stdin)
			} else {
				data, err = ioutil.ReadFile(rulesFile)
			}
			if err != nil {
				return nil, nil, err
			}
			var rules []schema.FirewallRule
			err = json.Unmarshal(data, &rules)
			if err != nil {
				return nil, nil, err
			}
			for _, rule := range rules {
				var sourceNets []net.IPNet
				for i, sourceIP := range rule.SourceIPs {
					_, sourceNet, err := net.ParseCIDR(sourceIP)
					if err != nil {
						return nil, nil, fmt.Errorf("invalid CIDR on index %d : %s", i, err)
					}
					sourceNets = append(sourceNets, *sourceNet)
				}
				opts.Rules = append(opts.Rules, hcloud.FirewallRule{
					Direction:   hcloud.FirewallRuleDirection(rule.Direction),
					SourceIPs:   sourceNets,
					Protocol:    hcloud.FirewallRuleProtocol(rule.Protocol),
					Port:        rule.Port,
					Description: rule.Description,
				})
			}
		}

		result, _, err := s.Client().Firewall().Create(s, opts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(cmd, s, result.Actions); err != nil {
			return nil, nil, err
		}

		cmd.Printf("Firewall %d created\n", result.Firewall.ID)

		return result.Firewall, util.Wrap("firewall", hcloud.SchemaFromFirewall(result.Firewall)), err
	},
}
