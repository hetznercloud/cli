package cmds

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var imageListTableOutput *output.Table
var typeFilter string

func init() {
	imageListTableOutput = describeImageListTableOutput(nil)
}

func newImageListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List images",
		Long: util.ListLongDescription(
			"Displays a list of images.",
			imageListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runImageList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(imageListTableOutput.Columns()), output.OptionJSON())
	cmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Only show images of given type")
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runImageList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.ImageListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
		IncludeDeprecated: true,
	}
	images, err := cli.Client().Image.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}
	if typeFilter != "" {
		var _images []*hcloud.Image
		for _, image := range images {
			if string(image.Type) == typeFilter {
				_images = append(_images, image)
			}
		}
		images = _images
	}

	if outOpts.IsSet("json") {
		var imageSchemas []schema.Image
		for _, image := range images {
			imageSchemas = append(imageSchemas, util.ImageToSchema(*image))
		}
		return util.DescribeJSON(imageSchemas)
	}

	cols := []string{"id", "type", "name", "description", "image_size", "disk_size", "created", "deprecated"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := describeImageListTableOutput(cli)
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, image := range images {
		tw.Write(cols, image)
	}
	tw.Flush()

	return nil
}

func describeImageListTableOutput(cli *state.State) *output.Table {
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
			if image.BoundTo != nil && cli != nil {
				return cli.ServerName(image.BoundTo.ID)
			}
			return util.NA("")
		})).
		AddFieldFn("created_from", output.FieldFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.CreatedFrom != nil && cli != nil {
				return cli.ServerName(image.CreatedFrom.ID)
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
		AddFieldFn("deprecated", output.FieldFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.Deprecated.IsZero() {
				return "-"
			}
			return util.Datetime(image.Deprecated)
		}))
}
