package network

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
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.Network, schema.Network]{
	ResourceNamePlural: "Networks",
	JSONKeyGetByName:   "networks",
	DefaultColumns:     []string{"id", "name", "ip_range", "servers", "age"},
	SortOption:         nil, // Networks do not support sorting

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.Network, error) {
		opts := hcloud.NetworkListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().Network().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.Network], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.Network{}).
			AddFieldFn("servers", func(network *hcloud.Network) string {
				serverCount := len(network.Servers)
				if serverCount == 1 {
					return fmt.Sprintf("%v server", serverCount)
				}
				return fmt.Sprintf("%v servers", serverCount)
			}).
			AddFieldFn("ip_range", func(network *hcloud.Network) string {
				return network.IPRange.String()
			}).
			AddFieldFn("labels", func(network *hcloud.Network) string {
				return util.LabelsToString(network.Labels)
			}).
			AddFieldFn("protection", func(network *hcloud.Network) string {
				var protection []string
				if network.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			}).
			AddFieldFn("created", func(network *hcloud.Network) string {
				return util.Datetime(network.Created)
			}).
			AddFieldFn("age", func(network *hcloud.Network) string {
				return util.Age(network.Created, time.Now())
			})
	},
	Schema: hcloud.SchemaFromNetwork,
}
