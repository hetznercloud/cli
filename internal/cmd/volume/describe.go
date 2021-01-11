package volume

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] VOLUME",
		Short:                 "Describe a volume",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.VolumeNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runVolumeDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runVolumeDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	volume, resp, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	switch {
	case outputFlags.IsSet("json"):
		return volumeDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(volume, outputFlags["format"][0])
	default:
		return volumeDescribeText(cli, volume)
	}
}

func volumeDescribeText(cli *state.State, volume *hcloud.Volume) error {
	fmt.Printf("ID:\t\t%d\n", volume.ID)
	fmt.Printf("Name:\t\t%s\n", volume.Name)
	fmt.Printf("Created:\t%s (%s)\n", util.Datetime(volume.Created), humanize.Time(volume.Created))
	fmt.Printf("Size:\t\t%s\n", humanize.Bytes(uint64(volume.Size*humanize.GByte)))
	fmt.Printf("Linux Device:\t%s\n", volume.LinuxDevice)
	fmt.Printf("Location:\n")
	fmt.Printf("  Name:\t\t%s\n", volume.Location.Name)
	fmt.Printf("  Description:\t%s\n", volume.Location.Description)
	fmt.Printf("  Country:\t%s\n", volume.Location.Country)
	fmt.Printf("  City:\t\t%s\n", volume.Location.City)
	fmt.Printf("  Latitude:\t%f\n", volume.Location.Latitude)
	fmt.Printf("  Longitude:\t%f\n", volume.Location.Longitude)
	if volume.Server != nil {
		server, _, err := cli.Client().Server.GetByID(cli.Context, volume.Server.ID)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %d", volume.Server.ID)
		}
		fmt.Printf("Server:\n")
		fmt.Printf("  ID:\t\t%d\n", server.ID)
		fmt.Printf("  Name:\t\t%s\n", server.Name)
	} else {
		fmt.Print("Server:\n  Not attached\n")
	}
	fmt.Printf("Protection:\n")
	fmt.Printf("  Delete:\t%s\n", util.YesNo(volume.Protection.Delete))

	fmt.Print("Labels:\n")
	if len(volume.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range volume.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	return nil
}

func volumeDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if volume, ok := data["volume"]; ok {
		return util.DescribeJSON(volume)
	}
	if volumes, ok := data["volumes"].([]interface{}); ok {
		return util.DescribeJSON(volumes[0])
	}
	return util.DescribeJSON(data)
}
