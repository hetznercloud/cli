package floatingip

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

var ListCmd = &base.ListCmd[*hcloud.FloatingIP, schema.FloatingIP]{
	ResourceNamePlural: "Floating IPs",
	JSONKeyGetByName:   "floating_ips",
	DefaultColumns:     []string{"id", "type", "name", "description", "ip", "home", "server", "dns", "age"},
	SortOption:         config.OptionSortFloatingIP,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.FloatingIP, error) {
		opts := hcloud.FloatingIPListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().FloatingIP().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.FloatingIP], client hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.FloatingIP{}).
			AddFieldFn("dns", func(floatingIP *hcloud.FloatingIP) string {
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
			}).
			AddFieldFn("server", func(floatingIP *hcloud.FloatingIP) string {
				var server string
				if floatingIP.Server != nil {
					return client.Server().ServerName(floatingIP.Server.ID)
				}
				return util.NA(server)
			}).
			AddFieldFn("home", func(floatingIP *hcloud.FloatingIP) string {
				return floatingIP.HomeLocation.Name
			}).
			AddFieldFn("ip", func(floatingIP *hcloud.FloatingIP) string {
				// Format IPv6 correctly
				if floatingIP.Network != nil {
					return floatingIP.Network.String()
				}
				return floatingIP.IP.String()
			}).
			AddFieldFn("protection", func(floatingIP *hcloud.FloatingIP) string {
				var protection []string
				if floatingIP.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			}).
			AddFieldFn("labels", func(floatingIP *hcloud.FloatingIP) string {
				return util.LabelsToString(floatingIP.Labels)
			}).
			AddFieldFn("created", func(floatingIP *hcloud.FloatingIP) string {
				return util.Datetime(floatingIP.Created)
			}).
			AddFieldFn("age", func(floatingIP *hcloud.FloatingIP) string {
				return util.Age(floatingIP.Created, time.Now())
			})
	},

	Schema: hcloud.SchemaFromFloatingIP,
}
