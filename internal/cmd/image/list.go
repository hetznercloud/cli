package image

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "images",
	DefaultColumns:     []string{"id", "type", "name", "description", "architecture", "image_size", "disk_size", "created", "deprecated"},
	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringP("type", "t", "", "Only show images of given type")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("backup", "snapshot", "system", "app"))

		cmd.Flags().StringSliceP("architecture", "a", []string{}, "Only show images of given architecture: x86|arm")
		cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))
	},
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.ImageListOpts{ListOpts: listOpts, IncludeDeprecated: true}

		imageType, _ := cmd.Flags().GetString("type")
		if len(imageType) > 0 {
			opts.Type = []hcloud.ImageType{hcloud.ImageType(imageType)}
		}

		architecture, _ := cmd.Flags().GetStringSlice("architecture")
		if len(architecture) > 0 {
			for _, arch := range architecture {
				opts.Architecture = append(opts.Architecture, hcloud.Architecture(arch))
			}
		}

		if len(sorts) > 0 {
			opts.Sort = sorts
		}

		images, err := client.Image().AllWithOpts(ctx, opts)

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

	JSONSchema: func(resources []interface{}) interface{} {
		var imageSchemas []schema.Image
		for _, resource := range resources {
			image := resource.(*hcloud.Image)
			imageSchemas = append(imageSchemas, util.ImageToSchema(*image))
		}
		return imageSchemas
	},
}
