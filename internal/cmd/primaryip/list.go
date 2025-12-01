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

var ListCmd = &base.ListCmd[*hcloud.PrimaryIP, schema.PrimaryIP]{
	ResourceNamePlural: "Primary IPs",
	JSONKeyGetByName:   "primary_ips",
	DefaultColumns:     []string{"id", "type", "name", "ip", "assignee", "auto_delete", "age"},
	SortOption:         config.OptionSortPrimaryIP,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.PrimaryIP, error) {
		opts := hcloud.PrimaryIPListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().PrimaryIP().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.PrimaryIP], client hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.PrimaryIP{}).
			AddFieldFn("ip", func(primaryIP *hcloud.PrimaryIP) string {
				// Format IPv6 correctly
				if primaryIP.Network != nil {
					return primaryIP.Network.String()
				}
				return primaryIP.IP.String()
			}).
			AddFieldFn("dns", func(primaryIP *hcloud.PrimaryIP) string {
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
			}).
			AddFieldFn("assignee", func(primaryIP *hcloud.PrimaryIP) string {
				assignee := ""
				if primaryIP.AssigneeID != 0 {
					switch primaryIP.AssigneeType {
					case "server":
						assignee = fmt.Sprintf("Server %s", client.Server().ServerName(primaryIP.AssigneeID))
					}
				}
				return util.NA(assignee)
			}).
			AddFieldFn("protection", func(primaryIP *hcloud.PrimaryIP) string {
				var protection []string
				if primaryIP.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			}).
			AddFieldFn("auto_delete", func(primaryIP *hcloud.PrimaryIP) string {
				return util.YesNo(primaryIP.AutoDelete)
			}).
			AddFieldFn("labels", func(primaryIP *hcloud.PrimaryIP) string {
				return util.LabelsToString(primaryIP.Labels)
			}).
			AddFieldFn("created", func(primaryIP *hcloud.PrimaryIP) string {
				return util.Datetime(primaryIP.Created)
			}).
			AddFieldFn("age", func(primaryIP *hcloud.PrimaryIP) string {
				return util.Age(primaryIP.Created, time.Now())
			})
	},

	Schema: hcloud.SchemaFromPrimaryIP,
}
