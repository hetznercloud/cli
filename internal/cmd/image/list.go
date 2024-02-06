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
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Images",
	JSONKeyGetByName:   "images",
	DefaultColumns:     []string{"id", "type", "name", "description", "architecture", "image_size", "disk_size", "created", "deprecated"},
	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringSliceP("type", "t", []string{}, "Only show images of given type: system|app|snapshot|backup")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("backup", "snapshot", "system", "app"))

		cmd.Flags().StringSliceP("architecture", "a", []string{}, "Only show images of given architecture: x86|arm")
		cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))
	},
	Fetch: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
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
			return nil, fmt.Errorf("unknown image type: %s\n", strings.Join(unknown, ", "))
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

		images, err := s.Client().Image().AllWithOpts(s, opts)

		var resources []interface{}
		for _, n := range images {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Image{}).
			AddFieldAlias("imagesize", "image size").
			AddFieldAlias("disksize", "disk size").
			AddFieldAlias("osflavor", "os flavor").
			AddFieldAlias("osversion", "os version").
			AddFieldAlias("rapiddeploy", "rapid deploy").
			AddFieldAlias("createdfrom", "created from").
			AddFieldAlias("boundto", "bound to").
			AddFieldFn("name", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				return util.NA(image.Name)
			})).
			AddFieldFn("image_size", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				if image.ImageSize == 0 {
					return util.NA("")
				}
				return fmt.Sprintf("%.2f GB", image.ImageSize)
			})).
			AddFieldFn("disk_size", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				return fmt.Sprintf("%.0f GB", image.DiskSize)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				return humanize.Time(image.Created)
			})).
			AddFieldFn("bound_to", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				if image.BoundTo != nil {
					return client.Server().ServerName(image.BoundTo.ID)
				}
				return util.NA("")
			})).
			AddFieldFn("created_from", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				if image.CreatedFrom != nil {
					return client.Server().ServerName(image.CreatedFrom.ID)
				}
				return util.NA("")
			})).
			AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				var protection []string
				if image.Protection.Delete {
					protection = append(protection, "delete")
				}
				return strings.Join(protection, ", ")
			})).
			AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				return util.LabelsToString(image.Labels)
			})).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				return util.Datetime(image.Created)
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				return util.Age(image.Created, time.Now())
			})).
			AddFieldFn("deprecated", output.FieldFn(func(obj interface{}) string {
				image := obj.(*hcloud.Image)
				if image.Deprecated.IsZero() {
					return "-"
				}
				return util.Datetime(image.Deprecated)
			}))
	},

	Schema: func(resources []interface{}) interface{} {
		imageSchemas := make([]schema.Image, 0, len(resources))
		for _, resource := range resources {
			image := resource.(*hcloud.Image)
			imageSchemas = append(imageSchemas, hcloud.SchemaFromImage(image))
		}
		return imageSchemas
	},
}
