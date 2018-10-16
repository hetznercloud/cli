package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVolumeDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] VOLUME",
		Short:                 "Delete a volume",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runVolumeDelete),
	}
	return cmd
}

func runVolumeDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	_, err = cli.Client().Volume.Delete(cli.Context, volume)
	if err != nil {
		return err
	}

	fmt.Printf("Volume %d deleted\n", volume.ID)
	return nil
}
