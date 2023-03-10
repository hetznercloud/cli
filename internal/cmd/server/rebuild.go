package server

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var RebuildCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "rebuild [FLAGS] SERVER",
			Short:                 "Rebuild a server",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("image", "", "ID or name of image to rebuild from (required)")
		cmd.RegisterFlagCompletionFunc("image", cmpl.SuggestCandidatesF(client.Image().Names))
		cmd.MarkFlagRequired("image")

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		serverIDOrName := args[0]
		server, _, err := client.Server().Get(ctx, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIDOrName)
		}

		imageIDOrName, _ := cmd.Flags().GetString("image")
		image, _, err := client.Image().Get(ctx, imageIDOrName)
		if err != nil {
			return err
		}
		if image == nil {
			return fmt.Errorf("image not found: %s", imageIDOrName)
		}

		opts := hcloud.ServerRebuildOpts{
			Image: image,
		}
		result, _, err := client.Server().RebuildWithResult(ctx, server, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, result.Action); err != nil {
			return err
		}

		fmt.Printf("Server %d rebuilt with image %s\n", server.ID, image.Name)

		// Only print the root password if it's not empty,
		// which is only the case if it wasn't created with an SSH key.
		if result.RootPassword != "" {
			fmt.Printf("Root password: %s\n", result.RootPassword)
		}

		return nil
	},
}
