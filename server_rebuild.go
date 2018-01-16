package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerRebuildCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rebuild [FLAGS] SERVER",
		Short:                 "Rebuild a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureActiveContext,
		RunE:                  cli.wrap(runServerRebuild),
	}

	cmd.Flags().String("image", "", "ID or name of image to rebuild from")
	cmd.Flag("image").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_image_names"},
	}
	cmd.MarkFlagRequired("image")

	return cmd
}

func runServerRebuild(cli *CLI, cmd *cobra.Command, args []string) error {
	serverIDOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, serverIDOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", serverIDOrName)
	}

	imageIDOrName, _ := cmd.Flags().GetString("image")
	image, _, err := cli.Client().Image.Get(cli.Context, imageIDOrName)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %s", imageIDOrName)
	}

	opts := hcloud.ServerRebuildOpts{
		Image: image,
	}
	action, _, err := cli.Client().Server.Rebuild(cli.Context, server, opts)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Server %s rebuilt with image %s\n", server.Name, imageIDOrName)
	return nil
}
