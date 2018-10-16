package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newVolumeDetachCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "detach [FLAGS] VOLUME",
		Short:                 "Detach a volume",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runVolumeDetach),
	}
	return cmd
}

func runVolumeDetach(cli *CLI, cmd *cobra.Command, args []string) error {
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	action, _, err := cli.Client().Volume.Detach(cli.Context, volume)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Volume %d detached\n", volume.ID)
	return nil
}
