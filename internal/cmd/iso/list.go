package iso

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.ISO, schema.ISO]{
	ResourceNamePlural: "ISOs",
	JSONKeyGetByName:   "isos",
	DefaultColumns:     []string{"id", "name", "description", "type", "architecture"},
	SortOption:         nil, // ISOs does not support sorting

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSlice("architecture", []string{}, "Only show Images of given architecture: x86|arm")
		_ = cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))

		cmd.Flags().Bool("include-architecture-wildcard", false, "Include ISOs with unknown architecture, only required if you want so show custom ISOs and still filter for architecture. (true, false)")

		cmd.Flags().StringSlice("type", []string{"public", "private"}, "Types to include (public, private)")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("public", "private"))
	},

	Fetch: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.ISO, error) {
		opts := hcloud.ISOListOpts{ListOpts: listOpts}

		types, _ := flags.GetStringSlice("type")

		var unknown []string
		for _, t := range types {
			switch t {
			case string(hcloud.ISOTypePublic), string(hcloud.ISOTypePrivate):
			default:
				unknown = append(unknown, t)
			}
		}
		if len(unknown) > 0 {
			return nil, fmt.Errorf("unknown ISO types %s", strings.Join(unknown, ", "))
		}

		architecture, _ := flags.GetStringSlice("architecture")
		if len(architecture) > 0 {
			for _, arch := range architecture {
				opts.Architecture = append(opts.Architecture, hcloud.Architecture(arch))
			}
		}

		includeArchitectureWildcard, _ := flags.GetBool("include-architecture-wildcard")
		if includeArchitectureWildcard {
			opts.IncludeArchitectureWildcard = includeArchitectureWildcard
		}

		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		return s.Client().ISO().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.ISO], _ hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.ISO{}).
			AddFieldFn("architecture", func(iso *hcloud.ISO) string {
				if iso.Architecture == nil {
					return "-"
				}
				return string(*iso.Architecture)
			})
	},

	Schema: hcloud.SchemaFromISO,
}
