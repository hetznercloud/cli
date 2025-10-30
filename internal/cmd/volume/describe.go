package volume

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/location"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Volume]{
	ResourceNameSingular: "Volume",
	ShortDescription:     "Describe a Volume",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Volume().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Volume, any, error) {
		v, _, err := s.Client().Volume().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return v, hcloud.SchemaFromVolume(v), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, volume *hcloud.Volume) error {

		fmt.Fprintf(out, "ID:\t%d\n", volume.ID)
		fmt.Fprintf(out, "Name:\t%s\n", volume.Name)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(volume.Created), humanize.Time(volume.Created))
		fmt.Fprintf(out, "Size:\t%s\n", humanize.Bytes(uint64(volume.Size)*humanize.GByte))
		fmt.Fprintf(out, "Linux Device:\t%s\n", volume.LinuxDevice)

		fmt.Fprintf(out, "Location:\t\n")
		fmt.Fprintf(out, "%s", util.PrefixLines(location.DescribeLocation(volume.Location), "  "))

		if volume.Server != nil {
			fmt.Fprintf(out, "Server:\t\n")
			fmt.Fprintf(out, "  ID:\t%d\n", volume.Server.ID)
			fmt.Fprintf(out, "  Name:\t%s\n", s.Client().Server().ServerName(volume.Server.ID))
		} else {
			fmt.Fprintf(out, "Server:\tNot attached\n")
		}

		fmt.Fprintf(out, "Protection:\t\n")
		fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(volume.Protection.Delete))

		util.DescribeLabels(out, volume.Labels, "")

		return nil
	},
}
