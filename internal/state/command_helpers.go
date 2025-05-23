package state

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:generate mockgen -package state -destination zz_command_helper_mock.go . ActionWaiter,TokenEnsurer

type ActionWaiter interface {
	WaitForActions(context.Context, *cobra.Command, ...*hcloud.Action) error
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
