package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerCreateImageCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create-image [FLAGS] SERVER",
		Short:                 "Create image from a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:    cli.wrap(runServerCreateImage),
		PreRunE: validateServerCreateImage,
	}
	cmd.Flags().String("type", "snapshot", "Image type")
	cmd.Flags().String("description", "", "Image description")
	return cmd
}

func validateServerCreateImage(cmd *cobra.Command, args []string) error {
	imageType, _ := cmd.Flags().GetString("type")
	switch imageType {
	case hcloud.ImageTypeBackup, hcloud.ImageTypeSnapshot:
		break
	default:
		return fmt.Errorf("invalid image type: %v", imageType)
	}

	return nil
}

func runServerCreateImage(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid server id")
	}

	imageType, _ := cmd.Flags().GetString("type")
	description, _ := cmd.Flags().GetString("description")

	server := &hcloud.Server{ID: id}
	opts := &hcloud.ServerCreateImageOpts{
		Type:        hcloud.ImageType(imageType),
		Description: hcloud.String(description),
	}
	result, _, err := cli.Client().Server.CreateImage(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}

	fmt.Printf("Image %d created from server %d\n", result.Image.ID, id)

	return nil
}
