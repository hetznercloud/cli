package zone

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.Zone, schema.Zone]{
	ResourceNamePlural: "Zones",
	JSONKeyGetByName:   "zones",
	DefaultColumns:     []string{"id", "name", "status", "mode", "record_count", "age"},
	SortOption:         config.OptionSortZone,

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("mode", "", "Only Zones with this mode are displayed")
		_ = cmd.RegisterFlagCompletionFunc("mode", cmpl.SuggestCandidates(string(hcloud.ZoneModePrimary), string(hcloud.ZoneModeSecondary)))
	},

	Fetch: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.Zone, error) {
		opts := hcloud.ZoneListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		if flags.Changed("mode") {
			mode, _ := flags.GetString("mode")

			if hcloud.ZoneMode(mode) != hcloud.ZoneModePrimary && hcloud.ZoneMode(mode) != hcloud.ZoneModeSecondary {
				return nil, fmt.Errorf("unknown Zone mode: %s", mode)
			}

			opts.Mode = hcloud.ZoneMode(mode)
		}

		return s.Client().Zone().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.Zone], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.Zone{}).
			AddFieldFn("name", func(zone *hcloud.Zone) string {
				return util.DisplayZoneName(zone.Name)
			}).
			AddFieldFn("name_idna", func(zone *hcloud.Zone) string {
				return zone.Name
			}).
			AddFieldFn("primary_nameservers", func(zone *hcloud.Zone) string {
				addressAndPorts := make([]string, 0, len(zone.PrimaryNameservers))
				for _, ns := range zone.PrimaryNameservers {

					addressAndPorts = append(addressAndPorts, net.JoinHostPort(ns.Address, strconv.Itoa(ns.Port)))
				}

				return strings.Join(addressAndPorts, ", ")
			}).
			AddFieldFn("authoritative_nameservers", func(zone *hcloud.Zone) string {
				return strings.Join(zone.AuthoritativeNameservers.Assigned, ", ")
			}).
			AddFieldFn("protection", func(zone *hcloud.Zone) string {
				var protection []string
				if zone.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			}).
			AddFieldFn("labels", func(zone *hcloud.Zone) string {
				return util.LabelsToString(zone.Labels)
			}).
			AddFieldFn("created", func(zone *hcloud.Zone) string {
				return util.Datetime(zone.Created)
			}).
			AddFieldFn("age", func(zone *hcloud.Zone) string {
				return util.Age(zone.Created, time.Now())
			})
	},

	Schema: hcloud.SchemaFromZone,
}
