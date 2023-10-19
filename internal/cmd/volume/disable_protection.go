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

var DisableProtectionCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "disable-protection [FLAGS] VOLUME PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Disable resource protection for a volume",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Volume().Names),
				cmpl.SuggestCandidates("delete"),
			),
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

		opts, err := getChangeProtectionOpts(false, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(ctx, client, waiter, volume, false, opts)
	},
}
