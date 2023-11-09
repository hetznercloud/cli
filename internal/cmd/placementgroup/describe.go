package placementgroup

import (
	"context"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "placement group",
	ShortDescription:     "Describe a placement group",
	JSONKeyGetByID:       "placement_group",
	JSONKeyGetByName:     "placement_groups",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PlacementGroup().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.PlacementGroup().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		placementGroup := resource.(*hcloud.PlacementGroup)

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
			cmd.Printf("    Server Name:\t%s\n", client.Server().ServerName(serverID))
		}

		cmd.Printf("Type:\t\t%s\n", placementGroup.Type)
		return nil
	},
}
