package datacenter

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Datacenter]{
	ResourceNameSingular: "Datacenter",
	ShortDescription:     "Describe a Datacenter",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Datacenter().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Datacenter, any, error) {
		dc, _, err := s.Client().Datacenter().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return dc, hcloud.SchemaFromDatacenter(dc), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, datacenter *hcloud.Datacenter) error {
		_, _ = fmt.Fprint(out, DescribeDatacenter(s.Client(), datacenter, false))
		return nil
	},
}

func DescribeDatacenter(client hcapi2.Client, datacenter *hcloud.Datacenter, short bool) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "ID:\t%d\n", datacenter.ID)
	_, _ = fmt.Fprintf(&sb, "Name:\t%s\n", datacenter.Name)
	_, _ = fmt.Fprintf(&sb, "Description:\t%s\n", datacenter.Description)
	_, _ = fmt.Fprintf(&sb, "Location:\t\n")
	_, _ = fmt.Fprintf(&sb, "  Name:\t%s\n", datacenter.Location.Name)
	_, _ = fmt.Fprintf(&sb, "  Description:\t%s\n", datacenter.Location.Description)
	_, _ = fmt.Fprintf(&sb, "  Country:\t%s\n", datacenter.Location.Country)
	_, _ = fmt.Fprintf(&sb, "  City:\t%s\n", datacenter.Location.City)
	_, _ = fmt.Fprintf(&sb, "  Latitude:\t%f\n", datacenter.Location.Latitude)
	_, _ = fmt.Fprintf(&sb, "  Longitude:\t%f\n", datacenter.Location.Longitude)

	if short {
		return sb.String()
	}

	type ServerTypeStatus struct {
		ID        int64
		Available bool
		Supported bool
	}

	allServerTypeStatus := make([]ServerTypeStatus, 0, len(datacenter.ServerTypes.Supported))
	for _, serverType := range datacenter.ServerTypes.Supported {
		allServerTypeStatus = append(allServerTypeStatus, ServerTypeStatus{ID: serverType.ID, Supported: true})
	}

	for _, serverType := range datacenter.ServerTypes.Available {
		index := slices.IndexFunc(allServerTypeStatus, func(i ServerTypeStatus) bool { return serverType.ID == i.ID })
		if index >= 0 {
			allServerTypeStatus[index].Available = true
		} else {
			allServerTypeStatus = append(allServerTypeStatus, ServerTypeStatus{ID: serverType.ID, Available: true})
		}
	}

	slices.SortFunc(allServerTypeStatus, func(a, b ServerTypeStatus) int { return int(a.ID - b.ID) })

	if len(allServerTypeStatus) > 0 {
		_, _ = fmt.Fprintf(&sb, "Server Types:\n")
		for _, t := range allServerTypeStatus {
			_, _ = fmt.Fprintf(&sb, "  - ID: %d\tName: %s\tSupported: %s\tAvailable: %s\n",
				t.ID,
				client.ServerType().ServerTypeName(t.ID),
				strconv.FormatBool(t.Supported),
				strconv.FormatBool(t.Available),
			)
		}
	} else {
		_, _ = fmt.Fprintf(&sb, "Server Types:\tNo Server Types\n")
	}

	return sb.String()
}
