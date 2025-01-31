package placementgroup

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.PlacementGroup]{
	ResourceNameSingular: "placement group",
	ShortDescription:     "Describe a placement group",
	JSONKeyGetByID:       "placement_group",
	JSONKeyGetByName:     "placement_groups",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.PlacementGroup, any, error) {
		pg, _, err := s.Client().PlacementGroup().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return pg, hcloud.SchemaFromPlacementGroup(pg), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, placementGroup *hcloud.PlacementGroup) error {
		cmd.Printf("ID:\t\t%d\n", placementGroup.ID)
		cmd.Printf("Name:\t\t%s\n", placementGroup.Name)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(placementGroup.Created), humanize.Time(placementGroup.Created))

		cmd.Print("Labels:\n")
		if len(placementGroup.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range placementGroup.Labels {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		cmd.Print("Servers:\n")
		for _, serverID := range placementGroup.Servers {
			cmd.Printf("  - Server ID:\t\t%d\n", serverID)
			cmd.Printf("    Server Name:\t%s\n", s.Client().Server().ServerName(serverID))
		}

		cmd.Printf("Type:\t\t%s\n", placementGroup.Type)
		return nil
	},
}
