package cli

import "github.com/spf13/cobra"

func newVolumeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "volume",
		Short:                 "Manage Volumes",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runVolume),
	}
	cmd.AddCommand(
		newVolumeListCommand(cli),
		newVolumeCreateCommand(cli),
		newVolumeUpdateCommand(cli),
		newVolumeDeleteCommand(cli),
		newVolumeDescribeCommand(cli),
		newVolumeAttachCommand(cli),
		newVolumeDetachCommand(cli),
		newVolumeResizeCommand(cli),
		newVolumeAddLabelCommand(cli),
		newVolumeRemoveLabelCommand(cli),
		newVolumeEnableProtectionCommand(cli),
		newVolumeDisableProtectionCommand(cli),
	)
	return cmd
}

func runVolume(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
