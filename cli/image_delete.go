package cli

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/spf13/cobra"
)

func newImageDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] IMAGE",
		Short:                 "Delete an image",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ImageNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runImageDelete),
	}
	return cmd
}

func runImageDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	image, _, err := cli.Client().Image.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %s", idOrName)
	}

	_, err = cli.Client().Image.Delete(cli.Context, image)
	if err != nil {
		return err
	}

	fmt.Printf("Image %d deleted\n", image.ID)
	return nil
}
