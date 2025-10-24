package image

import (
	"fmt"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
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

var ListCmd = &base.ListCmd[*hcloud.Image, schema.Image]{
	ResourceNamePlural: "Images",
	JSONKeyGetByName:   "images",
	DefaultColumns:     []string{"id", "type", "name", "description", "architecture", "image_size", "disk_size", "created", "deprecated"},
	SortOption:         config.OptionSortImage,

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSliceP("type", "t", []string{}, "Only show Images of given type: system|app|snapshot|backup")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("backup", "snapshot", "system", "app"))

		cmd.Flags().StringSliceP("architecture", "a", []string{}, "Only show Images of given architecture: x86|arm")
		_ = cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))
	},
	Fetch: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.Image, error) {
		opts := hcloud.ImageListOpts{ListOpts: listOpts, IncludeDeprecated: true}

		types, _ := flags.GetStringSlice("type")
		var (
			unknown []string
		)
		for _, imageType := range types {
			switch imageType {
			case string(hcloud.ImageTypeBackup), string(hcloud.ImageTypeSnapshot), string(hcloud.ImageTypeSystem), string(hcloud.ImageTypeApp):
				opts.Type = append(opts.Type, hcloud.ImageType(imageType))
			default:
				unknown = append(unknown, imageType)
			}
		}
		if len(unknown) > 0 {
			return nil, fmt.Errorf("unknown Image type: %s", strings.Join(unknown, ", "))
		}

		architecture, _ := flags.GetStringSlice("architecture")
		if len(architecture) > 0 {
			for _, arch := range architecture {
				opts.Architecture = append(opts.Architecture, hcloud.Architecture(arch))
			}
		}

		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		return s.Client().Image().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table[*hcloud.Image], client hcapi2.Client) {
		t.
			AddAllowedFields(&hcloud.Image{}).
			AddFieldAlias("imagesize", "image size").
			AddFieldAlias("disksize", "disk size").
			AddFieldAlias("osflavor", "os flavor").
			AddFieldAlias("osversion", "os version").
			AddFieldAlias("rapiddeploy", "rapid deploy").
			AddFieldAlias("createdfrom", "created from").
			AddFieldAlias("boundto", "bound to").
			AddFieldFn("name", func(image *hcloud.Image) string {
				return util.NA(image.Name)
			}).
			AddFieldFn("image_size", func(image *hcloud.Image) string {
				if image.ImageSize == 0 {
					return util.NA("")
				}
				return fmt.Sprintf("%.2f GB", image.ImageSize)
			}).
			AddFieldFn("disk_size", func(image *hcloud.Image) string {
				return fmt.Sprintf("%.0f GB", image.DiskSize)
			}).
			AddFieldFn("created", func(image *hcloud.Image) string {
				return humanize.Time(image.Created)
			}).
			AddFieldFn("bound_to", func(image *hcloud.Image) string {
				if image.BoundTo != nil {
					return client.Server().ServerName(image.BoundTo.ID)
				}
				return util.NA("")
			}).
			AddFieldFn("created_from", func(image *hcloud.Image) string {
				if image.CreatedFrom != nil {
					return client.Server().ServerName(image.CreatedFrom.ID)
				}
				return util.NA("")
			}).
			AddFieldFn("protection", func(image *hcloud.Image) string {
				var protection []string
				if image.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			}).
			AddFieldFn("labels", func(image *hcloud.Image) string {
				return util.LabelsToString(image.Labels)
			}).
			AddFieldFn("created", func(image *hcloud.Image) string {
				return util.Datetime(image.Created)
			}).
			AddFieldFn("age", func(image *hcloud.Image) string {
				return util.Age(image.Created, time.Now())
			}).
			AddFieldFn("deprecated", func(image *hcloud.Image) string {
				if image.Deprecated.IsZero() {
					return "-"
				}
				return util.Datetime(image.Deprecated)
			})
	},

	Schema: hcloud.SchemaFromImage,
}
