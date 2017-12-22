package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
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
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid image id")
	}
	image := &hcloud.Image{ID: id}

	_, err = cli.Client().Image.Delete(cli.Context, image)
	if err != nil {
		return err
	}

	fmt.Printf("Image %d deleted\n", id)
	return nil
}
