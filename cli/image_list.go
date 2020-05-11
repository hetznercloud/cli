package cli

import (
	"fmt"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var imageListTableOutput *tableOutput
var typeFilter string

func init() {
	imageListTableOutput = describeImageListTableOutput(nil)
}

func newImageListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List images",
		Long: listLongDescription(
			"Displays a list of images.",
			imageListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runImageList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(imageListTableOutput.Columns()), outputOptionJSON())
	cmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Only show images of given type")
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runImageList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.ImageListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
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
			imageSchemas = append(imageSchemas, imageToSchema(*image))
		}
		return describeJSON(imageSchemas)
	}

	cols := []string{"id", "type", "name", "description", "image_size", "disk_size", "created"}
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

func describeImageListTableOutput(cli *CLI) *tableOutput {
	return newTableOutput().
		AddAllowedFields(hcloud.Image{}).
		AddFieldAlias("imagesize", "image size").
		AddFieldAlias("disksize", "disk size").
		AddFieldAlias("osflavor", "os flavor").
		AddFieldAlias("osversion", "os version").
		AddFieldAlias("rapiddeploy", "rapid deploy").
		AddFieldAlias("createdfrom", "created from").
		AddFieldAlias("boundto", "bound to").
		AddFieldOutputFn("name", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			return na(image.Name)
		})).
		AddFieldOutputFn("image_size", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.ImageSize == 0 {
				return na("")
			}
			return fmt.Sprintf("%.1f GB", image.ImageSize)
		})).
		AddFieldOutputFn("disk_size", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			return fmt.Sprintf("%.0f GB", image.DiskSize)
		})).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			return humanize.Time(image.Created)
		})).
		AddFieldOutputFn("bound_to", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.BoundTo != nil && cli != nil {
				return cli.GetServerName(image.BoundTo.ID)
			}
			return na("")
		})).
		AddFieldOutputFn("created_from", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.CreatedFrom != nil && cli != nil {
				return cli.GetServerName(image.CreatedFrom.ID)
			}
			return na("")
		})).
		AddFieldOutputFn("protection", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			var protection []string
			if image.Protection.Delete {
				protection = append(protection, "delete")
			}
			return strings.Join(protection, ", ")
		})).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			return labelsToString(image.Labels)
		})).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			return datetime(image.Created)
		}))
}
