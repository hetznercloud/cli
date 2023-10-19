package firewall

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func newCreateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create FLAGS",
		Short:                 "Create a Firewall",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(cli.EnsureToken),
		RunE:                  cli.Wrap(runFirewallCreate),
	}
	cmd.Flags().String("name", "", "Name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	cmd.Flags().String("rules-file", "", "JSON file containing your routes (use - to read from stdin). The structure of the file needs to be the same as within the API: https://docs.hetzner.cloud/#firewalls-get-a-firewall ")
	return cmd
}

func runFirewallCreate(cli *state.State, cmd *cobra.Command, args []string) error {
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
			return err
		}
		var rules []schema.FirewallRule
		err = json.Unmarshal(data, &rules)
		if err != nil {
			return err
		}
		for _, rule := range rules {
			var sourceNets []net.IPNet
			for i, sourceIP := range rule.SourceIPs {
				_, sourceNet, err := net.ParseCIDR(sourceIP)
				if err != nil {
					return fmt.Errorf("invalid CIDR on index %d : %s", i, err)
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

	result, _, err := cli.Client().Firewall.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Firewall %d created\n", result.Firewall.ID)

	return nil
}
