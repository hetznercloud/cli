package floatingip

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "Floating IPs",
	DefaultColumns:     []string{"id", "type", "name", "description", "ip", "home", "server", "dns"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts) ([]interface{}, error) {
		floatingips, _, err := client.FloatingIP().List(ctx, hcloud.FloatingIPListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, n := range floatingips {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.FloatingIP{}).
			AddFieldFn("dns", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				var dns string
				if len(floatingIP.DNSPtr) == 1 {
					for _, v := range floatingIP.DNSPtr {
						dns = v
					}
				}
				if len(floatingIP.DNSPtr) > 1 {
					dns = fmt.Sprintf("%d entries", len(floatingIP.DNSPtr))
				}
				return util.NA(dns)
			})).
			AddFieldFn("server", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				var server string
				if floatingIP.Server != nil {
					return client.Server().ServerName(floatingIP.Server.ID)
				}
				return util.NA(server)
			})).
			AddFieldFn("home", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				return floatingIP.HomeLocation.Name
			})).
			AddFieldFn("ip", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				if floatingIP.Network != nil {
					return floatingIP.Network.String()
				}
				return floatingIP.IP.String()
			})).
			AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				var protection []string
				if floatingIP.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			})).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				return util.LabelsToString(floatingIP.Labels)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				return util.Datetime(floatingIP.Created)
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var floatingIPSchemas []schema.FloatingIP
		for _, resource := range resources {
			floatingIP := resource.(*hcloud.FloatingIP)
			floatingIPSchema := schema.FloatingIP{
				ID:           floatingIP.ID,
				Name:         floatingIP.Name,
				Description:  hcloud.String(floatingIP.Description),
				IP:           floatingIP.IP.String(),
				Created:      floatingIP.Created,
				Type:         string(floatingIP.Type),
				HomeLocation: util.LocationToSchema(*floatingIP.HomeLocation),
				Blocked:      floatingIP.Blocked,
				Protection:   schema.FloatingIPProtection{Delete: floatingIP.Protection.Delete},
				Labels:       floatingIP.Labels,
			}
			for ip, dnsPtr := range floatingIP.DNSPtr {
				floatingIPSchema.DNSPtr = append(floatingIPSchema.DNSPtr, schema.FloatingIPDNSPtr{
					IP:     ip,
					DNSPtr: dnsPtr,
				})
			}
			if floatingIP.Server != nil {
				floatingIPSchema.Server = hcloud.Int(floatingIP.Server.ID)
			}
			floatingIPSchemas = append(floatingIPSchemas, floatingIPSchema)
		}
		return floatingIPSchemas
	},
}
