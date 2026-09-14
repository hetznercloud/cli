package server

import (
	"fmt"

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
			Use:                   "rebuild [options] --image <image> <server>",
			Short:                 "Rebuild a server",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("image", "", "ID or name of Image to rebuild from (required)")
		_ = cmd.RegisterFlagCompletionFunc("image", cmpl.SuggestCandidatesF(client.Image().Names))
		_ = cmd.MarkFlagRequired("image")

		cmd.Flags().Bool("allow-deprecated-image", false, "Enable the use of deprecated images (default: false) (true, false)")

		cmd.Flags().StringArray("user-data-from-file", []string{}, "Read user data from specified file (use - to read from stdin)")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		serverIDOrName := args[0]
		server, _, err := s.Client().Server().Get(s, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", serverIDOrName)
		}

		imageIDOrName, _ := cmd.Flags().GetString("image")
		userDataFiles, _ := cmd.Flags().GetStringArray("user-data-from-file")

		// Select correct image based on Server Type architecture
		image, _, err := s.Client().Image().GetForArchitecture(s, imageIDOrName, server.ServerType.Architecture)
		if err != nil {
			return err
		}

		if image == nil {
			return fmt.Errorf("image %s for architecture %s not found", imageIDOrName, server.ServerType.Architecture)
		}

		allowDeprecatedImage, _ := cmd.Flags().GetBool("allow-deprecated-image")
		if warning, err := deprecatedImageWarning(image, allowDeprecatedImage); err != nil {
			return err
		} else if warning != "" {
			cmd.Print(warning)
		}

		opts := hcloud.ServerRebuildOpts{
			Image: image,
		}

		if len(userDataFiles) > 0 {
			userData, err := buildUserData(userDataFiles)
			if err != nil {
				return err
			}
			opts.UserData = &userData
		}

		result, _, err := s.Client().Server().RebuildWithResult(s, server, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return err
		}

		cmd.Printf("Server %d rebuilt with Image %s\n", server.ID, image.Name)

		// Only print the root password if it's not empty,
		// which is only the case if it wasn't created with an SSH key.
		if result.RootPassword != "" {
			cmd.Printf("Root password: %s\n", result.RootPassword)
		}

		return nil
	},
}
