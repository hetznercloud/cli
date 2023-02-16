package server

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

var AttachISOCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:              "attach-iso [FLAGS] SERVER ISO",
			Short:            "Attach an ISO to a server",
			Args:             cobra.ExactArgs(2),
			TraverseChildren: true,
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Server().Names),
				cmpl.SuggestCandidatesF(client.ISO().Names),
			),
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, command *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		isoIDOrName := args[1]
		iso, _, err := client.ISO().Get(ctx, isoIDOrName)
		if err != nil {
			return err
		}
		if iso == nil {
			return fmt.Errorf("ISO not found: %s", isoIDOrName)
		}

		action, _, err := client.Server().AttachISO(ctx, server, iso)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("ISO %s attached to server %d\n", iso.Name, server.ID)
		return nil
	},
}
