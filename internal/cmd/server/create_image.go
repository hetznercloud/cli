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

var CreateImageCmd = base.Cmd{
	BaseCobraCommand: func(hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create-image [options] --type <snapshot|backup> <server>",
			Short: "Create an image from a server",
		}
		cmd.Flags().String("type", "", "Image type (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("backup", "snapshot"))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("description", "", "Image description")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
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
		result, _, err := s.Client().Server().CreateImage(s, server, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return err
		}

		cmd.Printf("Image %d created from server %d\n", result.Image.ID, server.ID)

		return nil
	},
}
