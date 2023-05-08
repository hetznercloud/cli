package iso

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "isos",
	DefaultColumns:     []string{"id", "name", "description", "type", "architecture"},
	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSlice("architecture", []string{}, "Only show images of given architecture: x86|arm")
		cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))

		cmd.Flags().Bool("include-architecture-wildcard", false, "Include ISOs with unknown architecture, only required if you want so show custom ISOs and still filter for architecture.")
	},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.ISOListOpts{ListOpts: listOpts}

		architecture, _ := cmd.Flags().GetStringSlice("architecture")
		if len(architecture) > 0 {
			for _, arch := range architecture {
				opts.Architecture = append(opts.Architecture, hcloud.Architecture(arch))
			}
		}

		includeArchitectureWildcard, _ := cmd.Flags().GetBool("include-architecture-wildcard")
		if includeArchitectureWildcard {
			opts.IncludeWildcardArchitecture = includeArchitectureWildcard
		}

		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		isos, err := client.ISO().AllWithOpts(ctx, opts)

		var resources []interface{}
		for _, n := range isos {
			resources = append(resources, n)
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

	JSONSchema: func(resources []interface{}) interface{} {
		var isoSchemas []schema.ISO
		for _, resource := range resources {
			iso := resource.(*hcloud.ISO)
			isoSchemas = append(isoSchemas, util.ISOToSchema(*iso))
		}
		return isoSchemas
	},
}
