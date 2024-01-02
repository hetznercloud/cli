package primaryip

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var UnAssignCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "unassign [FLAGS] PRIMARYIP",
			Short: "Unassign a Primary IP from an assignee (usually a server)",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.PrimaryIP().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		primaryIP, _, err := client.PrimaryIP().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if primaryIP == nil {
			return fmt.Errorf("Primary IP not found: %v", idOrName)
		}

		action, _, err := client.PrimaryIP().Unassign(ctx, primaryIP.ID)
		if err != nil {
			return err
		}

		if err := actionWaiter.ActionProgress(cmd, ctx, action); err != nil {
			return err
		}

		cmd.Printf("Primary IP %d was unassigned successfully\n", primaryIP.ID)
		return nil
	},
}
