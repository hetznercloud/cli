package primaryip

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var AssignCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "assign [FLAGS] PRIMARYIP",
			Short: "Assign a Primary IP to an assignee (usually a server)",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.PrimaryIP().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("server", "", "Name or ID of the server")
		cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))
		cmd.MarkFlagRequired("server")
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

		serverIDOrName, _ := cmd.Flags().GetString("server")

		server, _, err := client.Server().Get(ctx, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIDOrName)
		}
		opts := hcloud.PrimaryIPAssignOpts{
			ID:           primaryIP.ID,
			AssigneeType: "server",
			AssigneeID:   server.ID,
		}

		action, _, err := client.PrimaryIP().Assign(ctx, opts)
		if err != nil {
			return err
		}

		if err := actionWaiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Primary IP %d assigned to %s %d\n", opts.ID, opts.AssigneeType, opts.AssigneeID)
		return nil
	},
}
