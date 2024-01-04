package volume

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "volume",
	ShortDescription:     "Describe an Volume",
	JSONKeyGetByID:       "volume",
	JSONKeyGetByName:     "volumes",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Volume().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		v, _, err := s.Client().Volume().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return v, hcloud.SchemaFromVolume(v), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		volume := resource.(*hcloud.Volume)

		cmd.Printf("ID:\t\t%d\n", volume.ID)
		cmd.Printf("Name:\t\t%s\n", volume.Name)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(volume.Created), humanize.Time(volume.Created))
		cmd.Printf("Size:\t\t%s\n", humanize.Bytes(uint64(volume.Size*humanize.GByte)))
		cmd.Printf("Linux Device:\t%s\n", volume.LinuxDevice)
		cmd.Printf("Location:\n")
		cmd.Printf("  Name:\t\t%s\n", volume.Location.Name)
		cmd.Printf("  Description:\t%s\n", volume.Location.Description)
		cmd.Printf("  Country:\t%s\n", volume.Location.Country)
		cmd.Printf("  City:\t\t%s\n", volume.Location.City)
		cmd.Printf("  Latitude:\t%f\n", volume.Location.Latitude)
		cmd.Printf("  Longitude:\t%f\n", volume.Location.Longitude)
		if volume.Server != nil {
			cmd.Printf("Server:\n")
			cmd.Printf("  ID:\t\t%d\n", volume.Server.ID)
			cmd.Printf("  Name:\t\t%s\n", s.Client().Server().ServerName(volume.Server.ID))
		} else {
			cmd.Print("Server:\n  Not attached\n")
		}
		cmd.Printf("Protection:\n")
		cmd.Printf("  Delete:\t%s\n", util.YesNo(volume.Protection.Delete))

		cmd.Print("Labels:\n")
		if len(volume.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range volume.Labels {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		return nil
	},
}
