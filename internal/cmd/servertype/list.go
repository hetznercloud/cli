package servertype

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.ServerType, schema.ServerType]{
	ResourceNamePlural: "Server Types",
	JSONKeyGetByName:   "server_types",
	DefaultColumns:     []string{"id", "name", "cores", "cpu_type", "architecture", "memory", "disk", "location"},
	SortOption:         nil, // Server Types do not support sorting

	Fetch: func(ctx context.Context, s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.ServerType, error) {
		opts := hcloud.ServerTypeListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().ServerType().AllWithOpts(ctx, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.ServerType], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.ServerType{}).
			AddFieldFn("location", func(serverType *hcloud.ServerType) string {
				now := time.Now()
				return listLocationNames(serverType, func(l hcloud.ServerTypeLocation) bool {
					return l.IsDeprecated() && l.UnavailableAfter().Before(now)
				})
			}).
			AddFieldFn("location_availability_signal", func(serverType *hcloud.ServerType) string {
				now := time.Now()
				return listLocationNames(serverType, func(l hcloud.ServerTypeLocation) bool {
					return (l.IsDeprecated() && l.UnavailableAfter().Before(now)) || !l.Available
				})
			}).
			AddFieldFn("location_recommended", func(serverType *hcloud.ServerType) string {
				return listLocationNames(serverType, func(l hcloud.ServerTypeLocation) bool {
					return !l.Recommended
				})
			}).
			AddFieldAlias("storagetype", "storage type").
			AddFieldFn("memory", func(serverType *hcloud.ServerType) string {
				return fmt.Sprintf("%.1f GB", serverType.Memory)
			}).
			AddFieldFn("disk", func(serverType *hcloud.ServerType) string {
				return fmt.Sprintf("%d GB", serverType.Disk)
			}).
			AddFieldFn("traffic", func(*hcloud.ServerType) string {
				// Was deprecated and traffic is now set per location, only available through describe.
				// Field was kept to avoid returning errors if people explicitly request the column.
				return "-"
			}).
			AddFieldFn("deprecated", func(serverType *hcloud.ServerType) string {
				deprecatedInfos := make([]string, 0, len(serverType.Locations))
				for _, loc := range serverType.Locations {
					if loc.IsDeprecated() {
						deprecatedInfos = append(
							deprecatedInfos,
							fmt.Sprintf("%s=%s", loc.Location.Name, loc.UnavailableAfter().Local().Format(time.DateOnly)),
						)
					}
				}

				if len(deprecatedInfos) > 0 {
					return strings.Join(deprecatedInfos, ",")
				}
				return "-"
			})
	},

	Schema: hcloud.SchemaFromServerType,
	Configure: func(_ state.State, cmd *cobra.Command) *cobra.Command {
		cmd.Long += `

The location column lists supported locations. location_availability_signal is an advisory snapshot of current
capacity only; it can change at any time and does not guarantee that a create request will succeed. The create
response is authoritative for the complete requested configuration.`
		return cmd
	},
}

func listLocationNames(serverType *hcloud.ServerType, del func(hcloud.ServerTypeLocation) bool) string {
	var locationNames []string
	for _, l := range serverType.Locations {
		if !del(l) {
			locationNames = append(locationNames, l.Location.Name)
		}
	}
	return strings.Join(locationNames, ", ")
}
