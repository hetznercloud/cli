package floatingip

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Floating IPs",
	JSONKeyGetByName:   "floating_ips",
	DefaultColumns:     []string{"id", "type", "name", "description", "ip", "home", "server", "dns", "age"},

	Fetch: func(ctx context.Context, client hcapi2.Client, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.FloatingIPListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		floatingIPs, err := client.FloatingIP().AllWithOpts(ctx, opts)

		var resources []interface{}
		for _, n := range floatingIPs {
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
				// Format IPv6 correctly
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
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				floatingIP := obj.(*hcloud.FloatingIP)
				return util.Age(floatingIP.Created, time.Now())
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		floatingIPSchemas := make([]schema.FloatingIP, 0, len(resources))
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
				floatingIPSchema.Server = hcloud.Ptr(floatingIP.Server.ID)
			}
			floatingIPSchemas = append(floatingIPSchemas, floatingIPSchema)
		}
		return floatingIPSchemas
	},
}
