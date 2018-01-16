package cli

import (
	"fmt"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newImageListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List images",
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runImageList),
	}
	return cmd
}

func runImageList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	images, err := cli.Client().Image.All(cli.Context)
	if err != nil {
		return err
	}

	cols := []string{"id", "type", "name", "description", "imagesize", "disksize", "created"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := newTableOutput().
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
		AddFieldOutputFn("image size", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.ImageSize == 0 {
				return na("")
			}
			return fmt.Sprintf("%.1f GB", image.ImageSize)
		})).
		AddFieldOutputFn("disk size", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			return fmt.Sprintf("%.0f GB", image.DiskSize)
		})).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			return humanize.Time(image.Created)
		})).
		AddFieldOutputFn("bound to", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.BoundTo != nil {
				return strconv.Itoa(image.BoundTo.ID)
			}
			return na("")
		})).
		AddFieldOutputFn("created from", fieldOutputFn(func(obj interface{}) string {
			image := obj.(*hcloud.Image)
			if image.CreatedFrom != nil {
				return strconv.Itoa(image.CreatedFrom.ID)
			}
			return na("")
		}))

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
