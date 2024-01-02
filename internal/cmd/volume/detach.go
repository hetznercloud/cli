package volume

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var DetachCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:                   "detach [FLAGS] VOLUME",
			Short:                 "Detach a volume",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Volume().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		volume, _, err := client.Volume().Get(ctx, args[0])
		if err != nil {
			return err
		}
		if volume == nil {
			return fmt.Errorf("volume not found: %s", args[0])
		}

		action, _, err := client.Volume().Detach(ctx, volume)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(cmd, ctx, action); err != nil {
			return err
		}

		cmd.Printf("Volume %d detached\n", volume.ID)
		return nil
	},
}
