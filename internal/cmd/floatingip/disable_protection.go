package floatingip

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
			Use:   "disable-protection [FLAGS] FLOATINGIP PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Disable resource protection for a Floating IP",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.FloatingIP().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		floatingIP, _, err := client.FloatingIP().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if floatingIP == nil {
			return fmt.Errorf("Floating IP not found: %v", idOrName)
		}

		opts, err := getChangeProtectionOpts(false, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(ctx, client, waiter, floatingIP, false, opts)
	},
}
