package config

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "Manage configuration",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		NewSetCommand(s),
		NewGetCommand(s),
		NewListCommand(s),
		NewUnsetCommand(s),
		NewAddCommand(s),
		NewRemoveCommand(s),
	)
	return cmd
}
