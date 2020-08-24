package cli

import "github.com/spf13/cobra"

func newImageCommand(cli *CLI) *cobra.Command {
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
