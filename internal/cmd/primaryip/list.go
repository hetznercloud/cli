package primaryip

import (
	"context"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "Primary IPs",
	DefaultColumns:     []string{"id", "type", "name", "ip", "assignee", "dns", "auto_delete"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.PrimaryIPListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		primaryips, _, err := client.PrimaryIP().List(ctx, opts)

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
			}))
	},
}
