package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newVolumeResizeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "resize [FLAGS] VOLUME",
		Short:                 "Resize a volume",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.VolumeNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runVolumeResize),
	}
	cmd.Flags().Int("size", 0, "New size (GB) of the volume (required)")
	cmd.MarkFlagRequired("size")
	return cmd
}

func runVolumeResize(cli *state.State, cmd *cobra.Command, args []string) error {
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
	fmt.Printf("You might need to adjust the filesystem size on the server too\n")
	return nil
}
