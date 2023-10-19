package base

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

// Cmd allows defining commands for generic resource-based commands
type Cmd struct {
	BaseCobraCommand func(hcapi2.Client) *cobra.Command
	Run              func(context.Context, hcapi2.Client, state.ActionWaiter, *cobra.Command, []string) error
}

// CobraCommand creates a command that can be registered with cobra.
func (gc *Cmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer, actionWaiter state.ActionWaiter,
) *cobra.Command {
	cmd := gc.BaseCobraCommand(client)

	if cmd.Args == nil {
		cmd.Args = cobra.NoArgs
	}

	cmd.TraverseChildren = true
	cmd.DisableFlagsInUseLine = true

	if cmd.PreRunE != nil {
		cmd.PreRunE = util.ChainRunE(cmd.PreRunE, tokenEnsurer.EnsureToken)
	} else {
		cmd.PreRunE = tokenEnsurer.EnsureToken
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return gc.Run(ctx, client, actionWaiter, cmd, args)
	}

	return cmd
}
