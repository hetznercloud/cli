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

var ReplaceRulesCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "replace-rules --rules-file <file> <firewall>",
			Short:                 "Replaces all rules from a Firewall from a file",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Firewall().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("rules-file", "", "JSON file containing your routes (use - to read from stdin). The structure of the file needs to be the same as within the API: https://docs.hetzner.cloud/#firewalls-get-a-firewall")
		_ = cmd.MarkFlagRequired("rules-file")
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

		opts := hcloud.FirewallSetRulesOpts{}
		rulesFile, _ := cmd.Flags().GetString("rules-file")
		if rulesFile != "" {
			rules, err := parseRulesFile(rulesFile)
			if err != nil {
				return err
			}
			opts.Rules = rules
		}

		actions, _, err := s.Client().Firewall().SetRules(s, firewall, opts)
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
