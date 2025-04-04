package server

import (
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
			Use:               "add-to-placement-group --placement-group <placement-group> <server>",
			Short:             "Add a Server to a Placement Group",
			ValidArgsFunction: cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
		}

		cmd.Flags().StringP("placement-group", "g", "", "Placement Group (ID or name) (required)")
		_ = cmd.RegisterFlagCompletionFunc("placement-group", cmpl.SuggestCandidatesF(client.PlacementGroup().Names))
		_ = cmd.MarkFlagRequired(("placement-group"))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found %s", idOrName)
		}

		placementGroupIDOrName, _ := cmd.Flags().GetString("placement-group")
		placementGroup, _, err := s.Client().PlacementGroup().Get(s, placementGroupIDOrName)
		if err != nil {
			return err
		}
		if placementGroup == nil {
			return fmt.Errorf("Placement Group not found: %s", placementGroupIDOrName)
		}

		action, _, err := s.Client().Server().AddToPlacementGroup(s, server, placementGroup)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Server %d added to Placement Group %s\n", server.ID, placementGroupIDOrName)
		return nil
	},
}
