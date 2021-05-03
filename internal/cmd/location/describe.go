package location

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// DescribeCmd defines a command for describing a location.
var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "location",
	ShortDescription:     "Describe a location",
	JSONKeyGetByID:       "location",
	JSONKeyGetByName:     "locations",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Location().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Location().Get(ctx, idOrName)
	},
	PrintText: func(ctx context.Context, client hcapi2.Client, resource interface{}) error {
		location := resource.(*hcloud.Location)

		fmt.Printf("ID:\t\t%d\n", location.ID)
		fmt.Printf("Name:\t\t%s\n", location.Name)
		fmt.Printf("Description:\t%s\n", location.Description)
		fmt.Printf("Network Zone:\t%s\n", location.NetworkZone)
		fmt.Printf("Country:\t%s\n", location.Country)
		fmt.Printf("City:\t\t%s\n", location.City)
		fmt.Printf("Latitude:\t%f\n", location.Latitude)
		fmt.Printf("Longitude:\t%f\n", location.Longitude)
		return nil
	},
}
