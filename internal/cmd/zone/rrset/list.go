package rrset

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.ZoneRRSet, schema.ZoneRRSet]{
	ResourceNamePlural: "Zone RRSets",
	JSONKeyGetByName:   "rrsets",

	DefaultColumns: []string{"name", "type", "records"},
	SortOption:     config.OptionSortZoneRRSet,

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSlice("type", nil, "Only Zone RRSets with one of these types are displayed")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates(rrsetTypeStrings...))
	},

	ValidArgsFunction: func(client hcapi2.Client) cobra.CompletionFunc {
		return cmpl.SuggestCandidatesF(client.Zone().Names)
	},

	PositionalArgumentOverride: []string{"zone"},

	FetchWithArgs: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string, args []string) ([]*hcloud.ZoneRRSet, error) {
		zoneIDOrName := args[0]
		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		types, _ := flags.GetStringSlice("type")

		zone, _, err := s.Client().Zone().Get(s, zoneIDOrName)
		if err != nil {
			return nil, err
		}
		if zone == nil {
			return nil, fmt.Errorf("Zone not found: %s", zoneIDOrName)
		}

		opts := hcloud.ZoneRRSetListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		if len(types) > 0 {
			for _, typ := range types {
				if !slices.Contains(rrsetTypeStrings, typ) {
					return nil, fmt.Errorf("unknown Zone RRSet type: %s", typ)
				}
				opts.Type = append(opts.Type, hcloud.ZoneRRSetType(typ))
			}
		}

		return s.Client().Zone().AllRRSetsWithOpts(s, zone, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.ZoneRRSet], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.ZoneRRSet{}).
			RemoveAllowedField("zone").
			AddFieldFn("protection", func(rrSet *hcloud.ZoneRRSet) string {
				var protection []string
				if rrSet.Protection.Change {
					protection = append(protection, "change")
				}
				return strings.Join(protection, ", ")
			}).
			AddFieldFn("labels", func(rrSet *hcloud.ZoneRRSet) string {
				return util.LabelsToString(rrSet.Labels)
			}).
			AddFieldFn("records", func(rrSet *hcloud.ZoneRRSet) string {
				var records []string
				for _, record := range rrSet.Records {
					records = append(records, record.Value)
				}
				return strings.Join(records, "\n")
			}).
			AddFieldFn("ttl", func(rrSet *hcloud.ZoneRRSet) string {
				if rrSet.TTL == nil {
					return "-"
				}
				return strconv.FormatInt(int64(*rrSet.TTL), 10)
			})
	},

	Schema: hcloud.SchemaFromZoneRRSet,

	Experimental: experimental.DNS,
}
