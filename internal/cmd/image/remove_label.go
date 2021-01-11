package image

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newRemoveLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label [FLAGS] IMAGE LABELKEY",
		Short: "Remove a label from an image",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.ImageNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.ImageLabelKeys(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateImageRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runImageRemoveLabel),
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

func runImageRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
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
