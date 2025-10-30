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
		_, _ = fmt.Fprint(out, DescribeLocation(location))
		return nil
	},
}

func DescribeLocation(location *hcloud.Location) string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "Location ID:\t%d\n", location.ID)
	_, _ = fmt.Fprintf(&sb, "Name:\t%s\n", location.Name)
	_, _ = fmt.Fprintf(&sb, "Description:\t%s\n", location.Description)
	_, _ = fmt.Fprintf(&sb, "Network Zone:\t%s\n", location.NetworkZone)
	_, _ = fmt.Fprintf(&sb, "Country:\t%s\n", location.Country)
	_, _ = fmt.Fprintf(&sb, "City:\t%s\n", location.City)
	_, _ = fmt.Fprintf(&sb, "Latitude:\t%f\n", location.Latitude)
	_, _ = fmt.Fprintf(&sb, "Longitude:\t%f\n", location.Longitude)
	return sb.String()
}
