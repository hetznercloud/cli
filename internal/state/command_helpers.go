package state

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

type ActionWaiter interface {
	ActionProgress(context.Context, *hcloud.Action) error
	WaitForActions(context.Context, []*hcloud.Action) error
}

type TokenEnsurer interface {
	EnsureToken(cmd *cobra.Command, args []string) error
}

func WrapCtx(
	ctx context.Context,
	fn func(context.Context, *cobra.Command, []string) error,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return fn(ctx, cmd, args)
	}
}
