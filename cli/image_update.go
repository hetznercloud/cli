package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newImageUpdateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] IMAGE",
		Short:                 "Update an image",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runImageUpdate),
	}

	cmd.Flags().String("description", "", "Image description")

	cmd.Flags().String("type", "", "Image type")
	cmd.Flag("type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_image_types_no_system"},
	}

	return cmd
}

func runImageUpdate(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	image, _, err := cli.Client().Image.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %s", idOrName)
	}

	description, _ := cmd.Flags().GetString("description")
	t, _ := cmd.Flags().GetString("type")
	opts := hcloud.ImageUpdateOpts{
		Description: hcloud.String(description),
		Type:        hcloud.ImageType(t),
	}
	_, _, err = cli.Client().Image.Update(cli.Context, image, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Image %d updated\n", image.ID)
	return nil
}
