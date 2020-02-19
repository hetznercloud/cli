package cli

import (
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var volumeListTableOutput *tableOutput

func init() {
	volumeListTableOutput = describeVolumeListTableOutput(nil)
}

func newVolumeListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List volumes",
		Long: listLongDescription(
			"Displays a list of volumes.",
			volumeListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runVolumeList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(volumeListTableOutput.Columns()), outputOptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runVolumeList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.VolumeListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	volumes, err := cli.Client().Volume.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var volumesSchema []schema.Volume
		for _, volume := range volumes {
			volumeSchema := schema.Volume{
				ID:          volume.ID,
				Name:        volume.Name,
				Location:    locationToSchema(*volume.Location),
				Size:        volume.Size,
				LinuxDevice: volume.LinuxDevice,
				Labels:      volume.Labels,
				Created:     volume.Created,
				Protection:  schema.VolumeProtection{Delete: volume.Protection.Delete},
			}
			if volume.Server != nil {
				volumeSchema.Server = hcloud.Int(volume.Server.ID)
			}
			volumesSchema = append(volumesSchema, volumeSchema)
		}
		return describeJSON(volumesSchema)
	}

	cols := []string{"id", "name", "size", "server", "location"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := describeVolumeListTableOutput(cli)
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, volume := range volumes {
		tw.Write(cols, volume)
	}
	tw.Flush()
	return nil
}

func describeVolumeListTableOutput(cli *CLI) *tableOutput {
	return newTableOutput().
		AddAllowedFields(hcloud.Volume{}).
		AddFieldOutputFn("server", fieldOutputFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			var server string
			if volume.Server != nil && cli != nil {
				return cli.GetServerName(volume.Server.ID)
			}
			return na(server)
		})).
		AddFieldOutputFn("size", fieldOutputFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			return humanize.Bytes(uint64(volume.Size * humanize.GByte))
		})).
		AddFieldOutputFn("location", fieldOutputFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			return volume.Location.Name
		})).
		AddFieldOutputFn("protection", fieldOutputFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			var protection []string
			if volume.Protection.Delete {
				protection = append(protection, "delete")
			}
			return strings.Join(protection, ", ")
		})).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			return labelsToString(volume.Labels)
		}))
}
