package primaryip

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Primary IPs",
	JSONKeyGetByName:   "primary_ips",
	DefaultColumns:     []string{"id", "type", "name", "ip", "assignee", "dns", "auto_delete", "age"},
	SortOption:         config.OptionSortPrimaryIP,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.PrimaryIPListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		primaryips, err := s.Client().PrimaryIP().AllWithOpts(s, opts)

		var resources []interface{}
		for _, n := range primaryips {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.PrimaryIP{}).
			AddFieldFn("ip", output.FieldFn(func(obj interface{}) string {
				primaryIP := obj.(*hcloud.PrimaryIP)
				// Format IPv6 correctly
				if primaryIP.Network != nil {
					return primaryIP.Network.String()
				}
				return primaryIP.IP.String()
			})).
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

	Schema: func(resources []interface{}) interface{} {
		primaryIPsSchema := make([]schema.PrimaryIP, 0, len(resources))
		for _, resource := range resources {
			primaryIP := resource.(*hcloud.PrimaryIP)
			primaryIPsSchema = append(primaryIPsSchema, hcloud.SchemaFromPrimaryIP(primaryIP))
		}
		return primaryIPsSchema
	},
}
