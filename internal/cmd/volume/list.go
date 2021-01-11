package volume

import (
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var volumeListTableOutput *output.Table

func init() {
	volumeListTableOutput = describeVolumeListTableOutput(nil)
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List volumes",
		Long: util.ListLongDescription(
			"Displays a list of volumes.",
			volumeListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runVolumeList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(volumeListTableOutput.Columns()), output.OptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runVolumeList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

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
				Location:    util.LocationToSchema(*volume.Location),
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
		return util.DescribeJSON(volumesSchema)
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

func describeVolumeListTableOutput(cli *state.State) *output.Table {
	return output.NewTable().
		AddAllowedFields(hcloud.Volume{}).
		AddFieldFn("server", output.FieldFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			var server string
			if volume.Server != nil && cli != nil {
				return cli.ServerName(volume.Server.ID)
			}
			return util.NA(server)
		})).
		AddFieldFn("size", output.FieldFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			return humanize.Bytes(uint64(volume.Size * humanize.GByte))
		})).
		AddFieldFn("location", output.FieldFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			return volume.Location.Name
		})).
		AddFieldFn("protection", output.FieldFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			var protection []string
			if volume.Protection.Delete {
				protection = append(protection, "delete")
			}
			return strings.Join(protection, ", ")
		})).
		AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			return util.LabelsToString(volume.Labels)
		})).
		AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			return util.Datetime(volume.Created)
		}))
}
