package volume

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "volume",
	ShortDescription:     "Describe an Volume",
	JSONKeyGetByID:       "volume",
	JSONKeyGetByName:     "volumes",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Volume().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Volume().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, client hcapi2.Client, resource interface{}) error {
		volume := resource.(*hcloud.Volume)

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
			fmt.Printf("Server:\n")
			fmt.Printf("  ID:\t\t%d\n", volume.Server.ID)
			fmt.Printf("  Name:\t\t%s\n", client.Server().ServerName(volume.Server.ID))
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
	},
}
