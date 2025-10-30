package location

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DescribeCmd defines a command for describing a location.
var DescribeCmd = base.DescribeCmd[*hcloud.Location]{
	ResourceNameSingular: "Location",
	ShortDescription:     "Describe a Location",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Location().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Location, any, error) {
		l, _, err := s.Client().Location().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return l, hcloud.SchemaFromLocation(l), nil
	},
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, location *hcloud.Location) error {
		fmt.Fprint(out, DescribeLocation(location))
		return nil
	},
}

func DescribeLocation(location *hcloud.Location) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "ID:\t%d\n", location.ID)
	fmt.Fprintf(&sb, "Name:\t%s\n", location.Name)
	fmt.Fprintf(&sb, "Description:\t%s\n", location.Description)
	fmt.Fprintf(&sb, "Network Zone:\t%s\n", location.NetworkZone)
	fmt.Fprintf(&sb, "Country:\t%s\n", location.Country)
	fmt.Fprintf(&sb, "City:\t%s\n", location.City)
	fmt.Fprintf(&sb, "Latitude:\t%f\n", location.Latitude)
	fmt.Fprintf(&sb, "Longitude:\t%f\n", location.Longitude)
	return sb.String()
}
