package server

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateImageCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create-image [FLAGS] SERVER",
			Short: "Create an image from a server",
			Args:  cobra.ExactArgs(1),
		}
		cmd.Flags().String("type", "", "Image type (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("backup", "snapshot"))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("description", "", "Image description")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		imageType, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		labels, _ := cmd.Flags().GetStringToString("label")

		switch hcloud.ImageType(imageType) {
		case hcloud.ImageTypeBackup, hcloud.ImageTypeSnapshot:
			break
		default:
			return fmt.Errorf("invalid image type: %v", imageType)
		}

		opts := &hcloud.ServerCreateImageOpts{
			Type:        hcloud.ImageType(imageType),
			Description: hcloud.Ptr(description),
			Labels:      labels,
		}
		result, _, err := client.Server().CreateImage(ctx, server, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, result.Action); err != nil {
			return err
		}

		fmt.Printf("Image %d created from server %d\n", result.Image.ID, server.ID)

		return nil
	},
}
