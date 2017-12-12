package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newImageDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "delete [flags] <id>",
		Short:            "Delete an image",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runImageDelete),
	}
	return cmd
}

func runImageDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid image id")
	}

	ctx := context.Background()
	_, err = cli.Client().Image.Delete(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("Image %d deleted\n", id)
	return nil
}
