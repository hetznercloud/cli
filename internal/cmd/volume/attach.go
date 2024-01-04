package volume

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var AttachCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "attach [FLAGS] VOLUME",
			Short:                 "Attach a volume to a server",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Volume().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("server", "", "Server (ID or name) (required)")
		cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))
		cmd.MarkFlagRequired("server")
		cmd.Flags().Bool("automount", false, "Automount volume after attach")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		volume, _, err := s.Client().Volume().Get(s, args[0])
		if err != nil {
			return err
		}
		if volume == nil {
			return fmt.Errorf("volume not found: %s", args[0])
		}

		serverIDOrName, _ := cmd.Flags().GetString("server")
		server, _, err := s.Client().Server().Get(s, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIDOrName)
		}
		automount, _ := cmd.Flags().GetBool("automount")
		action, _, err := s.Client().Volume().AttachWithOpts(s, volume, hcloud.VolumeAttachOpts{
			Server:    server,
			Automount: &automount,
		})

		if err != nil {
			return err
		}

		if err := s.ActionProgress(cmd, s, action); err != nil {
			return err
		}

		cmd.Printf("Volume %d attached to server %s\n", volume.ID, server.Name)
		return nil
	},
}
