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
		Args:                  util.ValidateExact,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newCreateCommand(s),
		newActiveCommand(s),
		newUseCommand(s),
		newDeleteCommand(s),
		newListCommand(s),
	)
	return cmd
}
