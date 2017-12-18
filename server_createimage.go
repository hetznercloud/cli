package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerCreateimageCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "create-image [flags] <id>",
		Short:            "Create image from a server",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runServerCreateimage),
		PreRunE:          validateServerCreateimage,
	}
	cmd.Flags().String("type", "snapshot", "Image type")
	cmd.Flags().String("description", "", "Image description")
	return cmd
}

func validateServerCreateimage(cmd *cobra.Command, args []string) error {
	imageType, _ := cmd.Flags().GetString("type")
	switch imageType {
	case hcloud.ImageTypeBackup, hcloud.ImageTypeSnapshot:
		break
	default:
		return fmt.Errorf("invalid image type: %v", imageType)
	}

	return nil
}

func runServerCreateimage(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	imageType, _ := cmd.Flags().GetString("type")
	description, _ := cmd.Flags().GetString("description")

	ctx := context.Background()
	server := &hcloud.Server{ID: id}
	opts := &hcloud.ServerCreateImageOpts{
		Type:        hcloud.ImageType(imageType),
		Description: hcloud.String(description),
	}
	result, _, err := cli.Client().Server.CreateImage(ctx, server, opts)
	if err != nil {
		return err
	}
	if err := <-waitAction(ctx, cli.Client(), result.Action); err != nil {
		return err
	}
	fmt.Printf("Image %d created from server %d\n", result.Image.ID, id)
	return nil
}
