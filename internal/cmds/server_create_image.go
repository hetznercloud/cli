package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerCreateImageCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create-image [FLAGS] SERVER",
		Short:                 "Create an image from a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateServerCreateImage, cli.EnsureToken),
		RunE:                  cli.Wrap(runServerCreateImage),
	}
	cmd.Flags().String("type", "", "Image type (required)")
	cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("backup", "snapshot"))
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("description", "", "Image description")

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	return cmd
}

func validateServerCreateImage(cmd *cobra.Command, args []string) error {
	imageType, _ := cmd.Flags().GetString("type")
	switch hcloud.ImageType(imageType) {
	case hcloud.ImageTypeBackup, hcloud.ImageTypeSnapshot:
		break
	default:
		return fmt.Errorf("invalid image type: %v", imageType)
	}

	return nil
}

func runServerCreateImage(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	imageType, _ := cmd.Flags().GetString("type")

	description, _ := cmd.Flags().GetString("description")

	labels, _ := cmd.Flags().GetStringToString("label")

	opts := &hcloud.ServerCreateImageOpts{
		Type:        hcloud.ImageType(imageType),
		Description: hcloud.String(description),
		Labels:      labels,
	}
	result, _, err := cli.Client().Server.CreateImage(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}

	fmt.Printf("Image %d created from server %d\n", result.Image.ID, server.ID)

	return nil
}
