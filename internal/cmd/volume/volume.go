package volume

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
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
