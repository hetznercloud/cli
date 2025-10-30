package placementgroup

import (
	"fmt"
	"io"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.PlacementGroup]{
	ResourceNameSingular: "Placement Group",
	ShortDescription:     "Describe a Placement Group",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.PlacementGroup, any, error) {
		pg, _, err := s.Client().PlacementGroup().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return pg, hcloud.SchemaFromPlacementGroup(pg), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, placementGroup *hcloud.PlacementGroup) error {
		_, _ = fmt.Fprintf(out, DescribePlacementGroup(s.Client(), placementGroup))
		return nil
	},
}

func DescribePlacementGroup(client hcapi2.Client, placementGroup *hcloud.PlacementGroup) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "ID:\t%d\n", placementGroup.ID)
	_, _ = fmt.Fprintf(&sb, "Name:\t%s\n", placementGroup.Name)
	_, _ = fmt.Fprintf(&sb, "Created:\t%s (%s)\n", util.Datetime(placementGroup.Created), humanize.Time(placementGroup.Created))
	_, _ = fmt.Fprintf(&sb, "Type:\t%s\n", placementGroup.Type)

	util.DescribeLabels(&sb, placementGroup.Labels, "")

	if len(placementGroup.Servers) == 0 {
		_, _ = fmt.Fprintf(&sb, "Servers:\tNo servers\n")
	} else {
		_, _ = fmt.Fprintf(&sb, "Servers:\t\n")
		for _, serverID := range placementGroup.Servers {
			_, _ = fmt.Fprintf(&sb, "  - Server ID:\t%d\n", serverID)
			_, _ = fmt.Fprintf(&sb, "    Server Name:\t%s\n", client.Server().ServerName(serverID))
		}
	}

	return sb.String()
}
