package cli

import (
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var volumeListTableOutput *tableOutput

func init() {
	volumeListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.Volume{}).
		AddFieldOutputFn("server", fieldOutputFn(func(obj interface{}) string {
			volume := obj.(*hcloud.Volume)
			var server string
			if volume.Server != nil {
				server = strconv.Itoa(volume.Server.ID)
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
	addListOutputFlag(cmd, volumeListTableOutput.Columns())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runVolumeList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.VolumeListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	sshKeys, err := cli.Client().Volume.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	cols := []string{"id", "name", "size", "server", "location"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := volumeListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, sshKey := range sshKeys {
		tw.Write(cols, sshKey)
	}
	tw.Flush()
	return nil
}
