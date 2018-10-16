package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newVolumeAddLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] VOLUME LABEL",
		Short:                 "Add a label to a volume",
		Args:                  cobra.ExactArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateVolumeAddLabel, cli.ensureToken),
		RunE:                  cli.wrap(runVolumeAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateVolumeAddLabel(cmd *cobra.Command, args []string) error {
	label := splitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runVolumeAddLabel(cli *CLI, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}
	label := splitLabel(args[1])

	if _, ok := volume.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on volume %d already exists", label[0], volume.ID)
	}
	labels := volume.Labels
	labels[label[0]] = label[1]
	opts := hcloud.VolumeUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Volume.Update(cli.Context, volume, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to volume %d\n", label[0], volume.ID)

	return nil
}
