package location

import (
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd[*hcloud.Location, schema.Location]{
	ResourceNamePlural: "Locations",
	JSONKeyGetByName:   "locations",
	DefaultColumns:     []string{"id", "name", "description", "network_zone", "country", "city"},
	SortOption:         config.OptionSortLocation,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.Location, error) {
		opts := hcloud.LocationListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().Location().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.Location{})
	},

	Schema: hcloud.SchemaFromLocation,
}
