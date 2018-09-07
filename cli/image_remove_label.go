package cli

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newImageRemoveLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-label [FLAGS] IMAGE LABELKEY",
		Short:                 "Remove a label from an image",
		Args:                  cobra.RangeArgs(1, 2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateImageRemoveLabel, cli.ensureToken),
		RunE:                  cli.wrap(runImageRemoveLabel),
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateImageRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	if all && len(args) == 2 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func runImageRemoveLabel(cli *CLI, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]
	image, _, err := cli.Client().Image.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %s", idOrName)
	}

	labels := image.Labels
	if all {
		labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := image.Labels[label]; !ok {
			return fmt.Errorf("label %s on image %d does not exist", label, image.ID)
		}
		delete(labels, label)
	}

	opts := hcloud.ImageUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Image.Update(cli.Context, image, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from image %d\n", image.ID)
	} else {
		fmt.Printf("Label %s removed from image %d\n", args[1], image.ID)
	}

	return nil
}
