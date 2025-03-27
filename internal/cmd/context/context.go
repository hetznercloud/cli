package context

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "context",
		Short:                 "Manage contexts",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		NewCreateCommand(s),
		NewActiveCommand(s),
		NewUnsetCommand(s),
		NewUseCommand(s),
		NewDeleteCommand(s),
		NewListCommand(s),
	)
	return cmd
}
