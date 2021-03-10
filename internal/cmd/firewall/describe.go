package firewall

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] FIREWALL",
		Short:                 "Describe a Firewall",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FirewallNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	idOrName := args[0]
	firewall, resp, err := cli.Client().Firewall.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if firewall == nil {
		return fmt.Errorf("Firewall not found: %v", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return describeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(firewall, outputFlags["format"][0])
	default:
		return describeText(cli, firewall)
	}
}

func describeText(cli *state.State, firewall *hcloud.Firewall) error {
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

	return nil
}

func describeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if firewall, ok := data["firewall"]; ok {
		return util.DescribeJSON(firewall)
	}
	if firewalls, ok := data["firewalls"].([]interface{}); ok {
		return util.DescribeJSON(firewalls[0])
	}
	return util.DescribeJSON(data)
}
