package servertype

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "server-type",
		Short:                 "Manage server types",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
	)
	return cmd
}
