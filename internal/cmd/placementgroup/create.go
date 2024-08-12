package placementgroup

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create [options] --name <name> --type <type>",
			Short: "Create a placement group",
		}
		cmd.Flags().String("name", "", "Name")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().String("type", "", "Type of the placement group")
		_ = cmd.MarkFlagRequired("type")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		labels, _ := cmd.Flags().GetStringToString("label")
		placementGroupType, _ := cmd.Flags().GetString("type")

		opts := hcloud.PlacementGroupCreateOpts{
			Name:   name,
			Labels: labels,
			Type:   hcloud.PlacementGroupType(placementGroupType),
		}

		result, _, err := s.Client().PlacementGroup().Create(s, opts)
		if err != nil {
			return nil, nil, err
		}

		if result.Action != nil {
			if err := s.WaitForActions(s, cmd, result.Action); err != nil {
				return nil, nil, err
			}
		}

		cmd.Printf("Placement group %d created\n", result.PlacementGroup.ID)

		return result.PlacementGroup, util.Wrap("placement_group", hcloud.SchemaFromPlacementGroup(result.PlacementGroup)), nil
	},
}
