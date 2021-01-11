package cmds

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newVolumeRemoveLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label [FLAGS] VOLUME LABELKEY",
		Short: "Remove a label from a volume",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.VolumeNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.VolumeLabelKeys(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateVolumeRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runVolumeRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateVolumeRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runVolumeRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	labels := volume.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := volume.Labels[label]; !ok {
			return fmt.Errorf("label %s on volume %d does not exist", label, volume.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.VolumeUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Volume.Update(cli.Context, volume, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from volume %d\n", volume.ID)
	} else {
		fmt.Printf("Label %s removed from volume %d\n", args[1], volume.ID)
	}

	return nil
}
