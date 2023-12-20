package placementgroup

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create FLAGS",
			Short: "Create a placement group",
		}
		cmd.Flags().String("name", "", "Name")
		cmd.MarkFlagRequired("name")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().String("type", "", "Type of the placement group")
		cmd.MarkFlagRequired("type")
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		labels, _ := cmd.Flags().GetStringToString("label")
		placementGroupType, _ := cmd.Flags().GetString("type")

		opts := hcloud.PlacementGroupCreateOpts{
			Name:   name,
			Labels: labels,
			Type:   hcloud.PlacementGroupType(placementGroupType),
		}

		res, _, err := client.PlacementGroup().Create(ctx, opts)
		if err != nil {
			return nil, nil, err
		}

		if res.Action != nil {
			if err := waiter.ActionProgress(ctx, res.Action); err != nil {
				return nil, nil, err
			}
		}

		cmd.Printf("Placement group %d created\n", res.PlacementGroup.ID)

		return res.PlacementGroup, util.Wrap("placement_group", hcloud.SchemaFromPlacementGroup(res.PlacementGroup)), nil
	},
}
