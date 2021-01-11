package volume

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newAttachCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "attach [FLAGS] VOLUME",
		Short:                 "Attach a volume to a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.VolumeNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runVolumeAttach),
	}
	cmd.Flags().String("server", "", "Server (ID or name) (required)")
	cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(cli.ServerNames))
	cmd.MarkFlagRequired("server")
	cmd.Flags().Bool("automount", false, "Automount volume after attach")

	return cmd
}

func runVolumeAttach(cli *state.State, cmd *cobra.Command, args []string) error {
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	serverIDOrName, _ := cmd.Flags().GetString("server")
	server, _, err := cli.Client().Server.Get(cli.Context, serverIDOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", serverIDOrName)
	}
	automount, _ := cmd.Flags().GetBool("automount")
	action, _, err := cli.Client().Volume.AttachWithOpts(cli.Context, volume, hcloud.VolumeAttachOpts{
		Server:    server,
		Automount: &automount,
	})

	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Volume %d attached to server %s\n", volume.ID, server.Name)
	return nil
}
