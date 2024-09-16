package firewall

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Firewalls",
	JSONKeyGetByName:   "firewalls",
	DefaultColumns:     []string{"id", "name", "rules_count", "applied_to_count"},
	SortOption:         config.OptionSortFirewall,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.FirewallListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		firewalls, err := s.Client().Firewall().AllWithOpts(s, opts)

		var resources []interface{}
		for _, n := range firewalls {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.Firewall{}).
			AddFieldFn("rules_count", output.FieldFn(func(obj interface{}) string {
				firewall := obj.(*hcloud.Firewall)
				count := len(firewall.Rules)
				if count == 1 {
					return fmt.Sprintf("%d Rule", count)
				}
				return fmt.Sprintf("%d Rules", count)
			})).
			AddFieldFn("applied_to_count", output.FieldFn(func(obj interface{}) string {
				firewall := obj.(*hcloud.Firewall)
				servers := 0
				labelSelectors := 0
				for _, r := range firewall.AppliedTo {
					if r.Type == hcloud.FirewallResourceTypeLabelSelector {
						labelSelectors++
						continue
					}
					servers++
				}
				serversText := "Servers"
				if servers == 1 {
					serversText = "Server"
				}
				labelSelectorsText := "Label Selectors"
				if labelSelectors == 1 {
					labelSelectorsText = "Label Selector"
				}
				return fmt.Sprintf("%d %s | %d %s", servers, serversText, labelSelectors, labelSelectorsText)
			}))
	},

	Schema: func(resources []interface{}) interface{} {
		firewallSchemas := make([]schema.Firewall, 0, len(resources))
		for _, resource := range resources {
			fw := resource.(*hcloud.Firewall)
			firewallSchemas = append(firewallSchemas, hcloud.SchemaFromFirewall(fw))
		}
		return firewallSchemas
	},
}
