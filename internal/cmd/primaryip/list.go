package primaryip

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
	ResourceNamePlural: "Primary IPs",
	DefaultColumns:     []string{"id", "type", "name", "ip", "assignee", "dns", "auto_delete", "age"},

	Fetch: func(ctx context.Context, client hcapi2.Client, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.PrimaryIPListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		primaryips, err := client.PrimaryIP().AllWithOpts(ctx, opts)

		var resources []interface{}
		for _, n := range primaryips {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.PrimaryIP{}).
			AddFieldFn("dns", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				var dns string
				if len(primaryIP.DNSPtr) == 1 {
					for _, v := range primaryIP.DNSPtr {
						dns = v
					}
				}
				if len(primaryIP.DNSPtr) > 1 {
					dns = fmt.Sprintf("%d entries", len(primaryIP.DNSPtr))
				}
				return util.NA(dns)
			})).
			AddFieldFn("assignee", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				assignee := ""
				if primaryIP.AssigneeID != 0 {
					switch primaryIP.AssigneeType {
					case "server":
						assignee = fmt.Sprintf("Server %s", client.Server().ServerName(primaryIP.AssigneeID))
					}
				}
				return util.NA(assignee)
			})).
			AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				var protection []string
				if primaryIP.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			})).
			AddFieldFn("auto_delete", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				return util.YesNo(primaryIP.AutoDelete)
			})).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				return util.LabelsToString(primaryIP.Labels)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				return util.Datetime(primaryIP.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				return util.Age(primaryIP.Created, time.Now())
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var primaryIPsSchema []schema.PrimaryIP
		for _, resource := range resources {
			primaryIP := resource.(*hcloud.PrimaryIP)
			var dnsPtrs []hcloud.PrimaryIPDNSPTR
			for i, d := range primaryIP.DNSPtr {
				dnsPtrs = append(dnsPtrs, hcloud.PrimaryIPDNSPTR{
					DNSPtr: d,
					IP:     i,
				})
			}
			var primaryIPSchema = schema.PrimaryIP{
				ID:           primaryIP.ID,
				Name:         primaryIP.Name,
				IP:           primaryIP.IP.String(),
				Type:         string(primaryIP.Type),
				AssigneeID:   primaryIP.AssigneeID,
				AssigneeType: primaryIP.AssigneeType,
				AutoDelete:   primaryIP.AutoDelete,
				Created:      primaryIP.Created,
				Datacenter:   util.DatacenterToSchema(*primaryIP.Datacenter),

				Protection: schema.PrimaryIPProtection{
					Delete: primaryIP.Protection.Delete,
				},
				Labels: primaryIP.Labels,
			}
			primaryIPsSchema = append(primaryIPsSchema, primaryIPSchema)
		}
		return primaryIPsSchema
	},
}
