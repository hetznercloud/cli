package primaryip

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "enable-protection PRIMARYIP",
			Short: "Enable Protection for a Primary IP",
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

		opts := hcloud.PrimaryIPChangeProtectionOpts{
			ID:     primaryIP.ID,
			Delete: true,
		}

		action, _, err := client.PrimaryIP().ChangeProtection(ctx, opts)
		if err != nil {
			return err
		}

		if err := actionWaiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Primary IP %d protection enabled", opts.ID)
		return nil
	},
}
