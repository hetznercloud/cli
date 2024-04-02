package server

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var RebuildCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "rebuild [--allow-deprecated-image] --image <image> <server>",
			Short:                 "Rebuild a server",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("image", "", "ID or name of image to rebuild from (required)")
		cmd.RegisterFlagCompletionFunc("image", cmpl.SuggestCandidatesF(client.Image().Names))
		cmd.MarkFlagRequired("image")
		cmd.Flags().Bool("allow-deprecated-image", false, "Enable the use of deprecated images (default: false)")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		serverIDOrName := args[0]
		server, _, err := s.Client().Server().Get(s, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIDOrName)
		}

		imageIDOrName, _ := cmd.Flags().GetString("image")
		// Select correct image based on server type architecture
		image, _, err := s.Client().Image().GetForArchitecture(s, imageIDOrName, server.ServerType.Architecture)
		if err != nil {
			return err
		}

		if image == nil {
			return fmt.Errorf("image %s for architecture %s not found", imageIDOrName, server.ServerType.Architecture)
		}

		allowDeprecatedImage, _ := cmd.Flags().GetBool("allow-deprecated-image")
		if !image.Deprecated.IsZero() {
			if allowDeprecatedImage {
				cmd.Printf("Attention: image %s is deprecated. It will continue to be available until %s.\n", image.Name, image.Deprecated.AddDate(0, 3, 0).Format(time.DateOnly))
			} else {
				return fmt.Errorf("image %s is deprecated, please use --allow-deprecated-image to create a server with this image. It will continue to be available until %s", image.Name, image.Deprecated.AddDate(0, 3, 0).Format(time.DateOnly))
			}
		}

		opts := hcloud.ServerRebuildOpts{
			Image: image,
		}
		result, _, err := s.Client().Server().RebuildWithResult(s, server, opts)
		if err != nil {
			return err
		}

		if err := s.ActionProgress(cmd, s, result.Action); err != nil {
			return err
		}

		cmd.Printf("Server %d rebuilt with image %s\n", server.ID, image.Name)

		// Only print the root password if it's not empty,
		// which is only the case if it wasn't created with an SSH key.
		if result.RootPassword != "" {
			cmd.Printf("Root password: %s\n", result.RootPassword)
		}

		return nil
	},
}
