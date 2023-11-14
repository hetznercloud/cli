package placementgroup

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
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
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) (*hcloud.Response, any, error) {
		name, _ := cmd.Flags().GetString("name")
		labels, _ := cmd.Flags().GetStringToString("label")
		placementGroupType, _ := cmd.Flags().GetString("type")

		opts := hcloud.PlacementGroupCreateOpts{
			Name:   name,
			Labels: labels,
			Type:   hcloud.PlacementGroupType(placementGroupType),
		}

		result, response, err := client.PlacementGroup().Create(ctx, opts)
		if err != nil {
			return nil, nil, err
		}

		if result.Action != nil {
			if err := waiter.ActionProgress(ctx, result.Action); err != nil {
				return nil, nil, err
			}
		}

		cmd.Printf("Placement group %d created\n", result.PlacementGroup.ID)

		return response, nil, nil
	},
	PrintResource: func(_ context.Context, _ hcapi2.Client, _ *cobra.Command, _ any) {
		// no-op
	},
}
