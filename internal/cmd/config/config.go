package config

import (
	"github.com/hetznercloud/cli/internal/cmd/config/sorting"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "config [FLAGS]",
		Short:                 "Manage config",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		sorting.NewSortCommand(cli),
	)
	return cmd
}
