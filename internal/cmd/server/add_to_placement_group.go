package server

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newAddToPlacementGroupCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-to-placement-group [FLAGS] SERVER",
		Short:                 "Add a server to a placement group",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runAddToPlacementGroup),
	}

	cmd.Flags().StringP("placement-group", "g", "", "Placement Group (ID or name) (required)")
	cmd.RegisterFlagCompletionFunc("placement-group", cmpl.SuggestCandidatesF(cli.PlacementGroupNames))
	cmd.MarkFlagRequired(("placement-group"))

	return cmd
}

func runAddToPlacementGroup(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found %s", idOrName)
	}

	placementGroupIDOrName, _ := cmd.Flags().GetString("placement_group")
	placementGroup, _, err := cli.Client().PlacementGroup.Get(cli.Context, placementGroupIDOrName)
	if err != nil {
		return err
	}
	if placementGroup == nil {
		return fmt.Errorf("placement group not found %s", placementGroupIDOrName)
	}

	action, _, err := cli.Client().Server.AddToPlacementGroup(cli.Context, server, placementGroup)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Server %d added to placement group %v", server.ID, placementGroupIDOrName)
	return nil
}
