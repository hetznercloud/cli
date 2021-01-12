package certificate

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "certificate",
		Short:                 "Manage certificates",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newListCommand(cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		newAddLabelCommand(cli),
		newRemoveLabelCommand(cli),
		newDeleteCommand(cli),
		newDescribeCommand(cli),
	)

	return cmd
}
