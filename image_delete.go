package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newImageDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] IMAGE",
		Short:                 "Delete an image",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runImageDelete),
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

	fmt.Printf("Image %s deleted\n", idOrName)
	return nil
}
