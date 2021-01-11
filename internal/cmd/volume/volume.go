package volume

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "volume",
		Short:                 "Manage Volumes",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newListCommand(cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		newDeleteCommand(cli),
		newDescribeCommand(cli),
		newAttachCommand(cli),
		newDetachCommand(cli),
		newResizeCommand(cli),
		newAddLabelCommand(cli),
		newRemoveLabelCommand(cli),
		newEnableProtectionCommand(cli),
		newDisableProtectionCommand(cli),
	)
	return cmd
}
