package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVolumeResizeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "resize [FLAGS] VOLUME",
		Short:                 "Resize a volume",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runVolumeResize),
	}
	cmd.Flags().Int("size", 0, "New size of the volume")
	cmd.MarkFlagRequired("size")
	return cmd
}

func runVolumeResize(cli *CLI, cmd *cobra.Command, args []string) error {
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	size, _ := cmd.Flags().GetInt("size")
	action, _, err := cli.Client().Volume.Resize(cli.Context, volume, size)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Volume %d resized\n", volume.ID)
	return nil
}
