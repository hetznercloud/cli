package firewall

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "Firewalls",
	DefaultColumns:     []string{"id", "name", "rules_count", "applied_to_count"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		firewalls, _, err := client.Firewall().List(ctx, hcloud.FirewallListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, n := range firewalls {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
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

	JSONSchema: func(resources []interface{}) interface{} {
		var firewallSchemas []schema.Firewall
		for _, resource := range resources {
			firewall := resource.(*hcloud.Firewall)
			firewallSchema := schema.Firewall{
				ID:      firewall.ID,
				Name:    firewall.Name,
				Labels:  firewall.Labels,
				Created: firewall.Created,
			}
			for _, rule := range firewall.Rules {
				var sourceNets []string
				for _, sourceIP := range rule.SourceIPs {
					sourceNets = append(sourceNets, sourceIP.Network())
				}
				var destinationNets []string
				for _, destinationIP := range rule.DestinationIPs {
					destinationNets = append(destinationNets, destinationIP.Network())
				}
				firewallSchema.Rules = append(firewallSchema.Rules, schema.FirewallRule{
					Direction:      string(rule.Direction),
					SourceIPs:      sourceNets,
					DestinationIPs: destinationNets,
					Protocol:       string(rule.Protocol),
					Port:           rule.Port,
				})
			}
			for _, AppliedTo := range firewall.AppliedTo {
				s := schema.FirewallResource{
					Type: string(AppliedTo.Type),
				}
				switch AppliedTo.Type {
				case hcloud.FirewallResourceTypeServer:
					s.Server = &schema.FirewallResourceServer{ID: AppliedTo.Server.ID}
				case hcloud.FirewallResourceTypeLabelSelector:
					s.LabelSelector = &schema.FirewallResourceLabelSelector{Selector: AppliedTo.LabelSelector.Selector}
				}

				firewallSchema.AppliedTo = append(firewallSchema.AppliedTo, s)
			}

			firewallSchemas = append(firewallSchemas, firewallSchema)
		}
		return firewallSchemas
	},
}
