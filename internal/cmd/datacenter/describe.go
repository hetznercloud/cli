package datacenter

import (
	"slices"
	"strconv"

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
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, interface{}, error) {
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

		type ServerTypeStatus struct {
			ID        int64
			Available bool
			Supported bool
		}

		allServerTypeStatus := make([]*ServerTypeStatus, 0, len(datacenter.ServerTypes.Supported))
		for _, serverType := range datacenter.ServerTypes.Supported {
			allServerTypeStatus = append(allServerTypeStatus, &ServerTypeStatus{ID: serverType.ID, Supported: true})
		}

		for _, serverType := range datacenter.ServerTypes.Available {
			index := slices.IndexFunc(allServerTypeStatus, func(i *ServerTypeStatus) bool { return serverType.ID == i.ID })
			if index >= 0 {
				allServerTypeStatus[index].Available = true
			} else {
				allServerTypeStatus = append(allServerTypeStatus, &ServerTypeStatus{ID: serverType.ID, Available: true})
			}
		}

		slices.SortFunc(allServerTypeStatus, func(a, b *ServerTypeStatus) int { return int(a.ID - b.ID) })

		if len(allServerTypeStatus) > 0 {
			cmd.Printf("Server Types:\n")
			for _, t := range allServerTypeStatus {
				cmd.Printf("  - ID: %-8d Name: %-8s Supported: %-8s Available: %s\n",
					t.ID,
					s.Client().ServerType().ServerTypeName(t.ID),
					strconv.FormatBool(t.Supported),
					strconv.FormatBool(t.Available),
				)
			}
		}

		return nil
	},
}
