package server

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var AddToPlacementGroupCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:               "add-to-placement-group [FLAGS] SERVER",
			Short:             "Add a server to a placement group",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
		}

		cmd.Flags().StringP("placement-group", "g", "", "Placement Group (ID or name) (required)")
		cmd.RegisterFlagCompletionFunc("placement-group", cmpl.SuggestCandidatesF(client.PlacementGroup().Names))
		cmd.MarkFlagRequired(("placement-group"))

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found %s", idOrName)
		}

		placementGroupIDOrName, _ := cmd.Flags().GetString("placement-group")
		placementGroup, _, err := client.PlacementGroup().Get(ctx, placementGroupIDOrName)
		if err != nil {
			return err
		}
		if placementGroup == nil {
			return fmt.Errorf("placement group not found %s", placementGroupIDOrName)
		}

		action, _, err := client.Server().AddToPlacementGroup(ctx, server, placementGroup)
		if err != nil {
			return err
		}

		if err := actionWaiter.ActionProgress(cmd, ctx, action); err != nil {
			return err
		}

		cmd.Printf("Server %d added to placement group %s\n", server.ID, placementGroupIDOrName)
		return nil
	},
}
