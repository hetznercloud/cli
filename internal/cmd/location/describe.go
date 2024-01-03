package location

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DescribeCmd defines a command for describing a location.
var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "location",
	ShortDescription:     "Describe a location",
	JSONKeyGetByID:       "location",
	JSONKeyGetByName:     "locations",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Location().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		l, _, err := s.Location().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return l, hcloud.SchemaFromLocation(l), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		location := resource.(*hcloud.Location)

		cmd.Printf("ID:\t\t%d\n", location.ID)
		cmd.Printf("Name:\t\t%s\n", location.Name)
		cmd.Printf("Description:\t%s\n", location.Description)
		cmd.Printf("Network Zone:\t%s\n", location.NetworkZone)
		cmd.Printf("Country:\t%s\n", location.Country)
		cmd.Printf("City:\t\t%s\n", location.City)
		cmd.Printf("Latitude:\t%f\n", location.Latitude)
		cmd.Printf("Longitude:\t%f\n", location.Longitude)
		return nil
	},
}
