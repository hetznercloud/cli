package placementgroup

import (
	"context"
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
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.PlacementGroup().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.PlacementGroup, any, error) {
		pg, _, err := s.Client().PlacementGroup().Get(cmd.Context(), idOrName)
		if err != nil {
			return nil, nil, err
		}
		return pg, hcloud.SchemaFromPlacementGroup(pg), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, out io.Writer, placementGroup *hcloud.PlacementGroup) error {
		description, err := DescribePlacementGroup(cmd.Context(), s.Client(), placementGroup)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(out, description)
		return err
	},
}

func DescribePlacementGroup(ctx context.Context, client hcapi2.Client, placementGroup *hcloud.PlacementGroup) (string, error) {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ID:\t%d\n", placementGroup.ID)
	fmt.Fprintf(&sb, "Name:\t%s\n", placementGroup.Name)
	fmt.Fprintf(&sb, "Created:\t%s (%s)\n", util.Datetime(placementGroup.Created), humanize.Time(placementGroup.Created))
	fmt.Fprintf(&sb, "Type:\t%s\n", placementGroup.Type)

	fmt.Fprintln(&sb)
	util.DescribeLabels(&sb, placementGroup.Labels, "")

	fmt.Fprintln(&sb)
	fmt.Fprintf(&sb, "Servers:\n")
	if len(placementGroup.Servers) == 0 {
		fmt.Fprintf(&sb, "  No servers\n")
	} else {
		for _, serverID := range placementGroup.Servers {
			name, err := client.Server().ServerName(ctx, serverID)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&sb, "  - Server ID:\t%d\n", serverID)
			fmt.Fprintf(&sb, "    Server Name:\t%s\n", name)
		}
	}

	return sb.String(), nil
}
