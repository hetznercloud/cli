package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewImageCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "image",
		Short:                 "Manage images",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newImageListCommand(cli),
		newImageDeleteCommand(cli),
		newImageDescribeCommand(cli),
		newImageUpdateCommand(cli),
		newImageEnableProtectionCommand(cli),
		newImageDisableProtectionCommand(cli),
		newImageAddLabelCommand(cli),
		newImageRemoveLabelCommand(cli),
	)
	return cmd
}
