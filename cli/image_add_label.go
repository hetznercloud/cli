package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newImageAddLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] IMAGE LABEL",
		Short:                 "Add a label to an image",
		Args:                  cobra.ExactArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateImageAddLabel, cli.ensureToken),
		RunE:                  cli.wrap(runImageAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateImageAddLabel(cmd *cobra.Command, args []string) error {
	label := splitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runImageAddLabel(cli *CLI, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]
	image, _, err := cli.Client().Image.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %s", idOrName)
	}
	label := splitLabel(args[1])

	if _, ok := image.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on image %d already exists", label[0], image.ID)
	}
	labels := image.Labels
	labels[label[0]] = label[1]
	opts := hcloud.ImageUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Image.Update(cli.Context, image, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to image %d\n", label[0], image.ID)

	return nil
}
