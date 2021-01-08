package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerRebuildCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rebuild [FLAGS] SERVER",
		Short:                 "Rebuild a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerRebuild),
	}

	cmd.Flags().String("image", "", "ID or name of image to rebuild from (required)")
	cmd.RegisterFlagCompletionFunc("image", cmpl.SuggestCandidatesF(cli.ImageNames))
	cmd.MarkFlagRequired("image")

	return cmd
}

func runServerRebuild(cli *state.State, cmd *cobra.Command, args []string) error {
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

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Server %d rebuilt with image %s\n", server.ID, image.Name)
	return nil
}
