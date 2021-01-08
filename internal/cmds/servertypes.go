package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewServerTypeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "server-type",
		Short:                 "Manage server types",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newServerTypeListCommand(cli),
		newServerTypeDescribeCommand(cli),
	)
	return cmd
}
