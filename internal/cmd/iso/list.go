package iso

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "ISOs",
	JSONKeyGetByName:   "isos",
	DefaultColumns:     []string{"id", "name", "description", "type", "architecture"},
	SortOption:         config.OptionSortISO,

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSlice("architecture", []string{}, "Only show images of given architecture: x86|arm")
		cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))

		cmd.Flags().Bool("include-architecture-wildcard", false, "Include ISOs with unknown architecture, only required if you want so show custom ISOs and still filter for architecture.")

		cmd.Flags().StringSlice("type", []string{"public", "private"}, "Types to include (public, private)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("public", "private"))
	},

	Fetch: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.ISOListOpts{ListOpts: listOpts}

		types, _ := flags.GetStringSlice("type")

		var unknown []string
		for _, t := range types {
			switch t {
			case string(hcloud.ISOTypePublic), string(hcloud.ISOTypePrivate):
				break
			default:
				unknown = append(unknown, t)
			}
		}
		if len(unknown) > 0 {
			return nil, fmt.Errorf("unknown ISO types %s\n", strings.Join(unknown, ", "))
		}

		architecture, _ := flags.GetStringSlice("architecture")
		if len(architecture) > 0 {
			for _, arch := range architecture {
				opts.Architecture = append(opts.Architecture, hcloud.Architecture(arch))
			}
		}

		includeArchitectureWildcard, _ := flags.GetBool("include-architecture-wildcard")
		if includeArchitectureWildcard {
			opts.IncludeWildcardArchitecture = includeArchitectureWildcard
		}

		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		isos, err := s.Client().ISO().AllWithOpts(s, opts)

		var resources []interface{}
		for _, iso := range isos {
			if slices.Contains(types, string(iso.Type)) {
				resources = append(resources, iso)
			}
		}
		return resources, err
	},

	OutputTable: func(_ hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.ISO{}).
			AddFieldFn("architecture", func(obj interface{}) string {
				iso := obj.(*hcloud.ISO)
				if iso.Architecture == nil {
					return "-"
				} else {
					return string(*iso.Architecture)
				}
			})
	},

	Schema: func(resources []interface{}) interface{} {
		isoSchemas := make([]schema.ISO, 0, len(resources))
		for _, resource := range resources {
			iso := resource.(*hcloud.ISO)
			isoSchemas = append(isoSchemas, hcloud.SchemaFromISO(iso))
		}
		return isoSchemas
	},
}
