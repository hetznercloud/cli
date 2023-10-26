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

var RemoveFromPlacementGroupCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:               "remove-from-placement-group SERVER",
			Short:             "Removes a server from a placement group",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
		}

		return cmd
	},

	Run: func(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		action, _, err := client.Server().RemoveFromPlacementGroup(ctx, server)
		if err != nil {
			return err
		}

		if err := actionWaiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Server %d removed from placement group\n", server.ID)
		return nil
	},
}
