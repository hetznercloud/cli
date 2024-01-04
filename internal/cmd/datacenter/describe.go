package datacenter

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "datacenter",
	ShortDescription:     "Describe an datacenter",
	JSONKeyGetByID:       "datacenter",
	JSONKeyGetByName:     "datacenters",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Datacenter().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		dc, _, err := s.Client().Datacenter().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return dc, hcloud.SchemaFromDatacenter(dc), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		datacenter := resource.(*hcloud.Datacenter)

		cmd.Printf("ID:\t\t%d\n", datacenter.ID)
		cmd.Printf("Name:\t\t%s\n", datacenter.Name)
		cmd.Printf("Description:\t%s\n", datacenter.Description)
		cmd.Printf("Location:\n")
		cmd.Printf("  Name:\t\t%s\n", datacenter.Location.Name)
		cmd.Printf("  Description:\t%s\n", datacenter.Location.Description)
		cmd.Printf("  Country:\t%s\n", datacenter.Location.Country)
		cmd.Printf("  City:\t\t%s\n", datacenter.Location.City)
		cmd.Printf("  Latitude:\t%f\n", datacenter.Location.Latitude)
		cmd.Printf("  Longitude:\t%f\n", datacenter.Location.Longitude)
		cmd.Printf("Server Types:\n")

		printServerTypes := func(list []*hcloud.ServerType) {
			for _, t := range list {
				cmd.Printf("  - ID:\t\t %d\n", t.ID)
				cmd.Printf("    Name:\t %s\n", s.Client().ServerType().ServerTypeName(t.ID))
				cmd.Printf("    Description: %s\n", s.Client().ServerType().ServerTypeDescription(t.ID))
			}
		}

		cmd.Printf("  Available:\n")
		if len(datacenter.ServerTypes.Available) > 0 {
			printServerTypes(datacenter.ServerTypes.Available)
		} else {
			cmd.Printf("    No available server types\n")
		}
		cmd.Printf("  Supported:\n")
		if len(datacenter.ServerTypes.Supported) > 0 {
			printServerTypes(datacenter.ServerTypes.Supported)
		} else {
			cmd.Printf("    No supported server types\n")
		}

		return nil
	},
}
