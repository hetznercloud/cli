package datacenter

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.Datacenter, schema.Datacenter]{
	ResourceNamePlural: "Datacenters",
	JSONKeyGetByName:   "datacenters",
	DefaultColumns:     []string{"id", "name", "description", "location"},
	SortOption:         config.OptionSortDatacenter,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.Datacenter, error) {
		opts := hcloud.DatacenterListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().Datacenter().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.Datacenter], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.Datacenter{}).
			AddFieldFn("location", func(datacenter *hcloud.Datacenter) string {
				return datacenter.Location.Name
			})
	},

	Schema: hcloud.SchemaFromDatacenter,

	Configure: func(s state.State, c *cobra.Command) *cobra.Command {
		c.Short += " (deprecated)"
		c.Long = `The 'hcloud datacenter ...' commands are deprecated and will be removed after 1 Oct. 2026.
After this date, requests to the datacenters API endpoints will return HTTP 410 Gone.

See https://docs.hetzner.cloud/changelog#2026-06-02-datacenters-deprecated for more details.

` + c.Long
		return c
	},
}
