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
			Use:                   "attach [--automount] --server <server> <volume>",
			Short:                 "Attach a Volume to a Server",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Volume().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("server", "", "Server (ID or name) (required)")
		_ = cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))
		_ = cmd.MarkFlagRequired("server")
		cmd.Flags().Bool("automount", false, "Automount Volume after attach (true, false)")

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
			return fmt.Errorf("Server not found: %s", serverIDOrName)
		}
		automount, _ := cmd.Flags().GetBool("automount")
		action, _, err := s.Client().Volume().AttachWithOpts(s, volume, hcloud.VolumeAttachOpts{
			Server:    server,
			Automount: &automount,
		})

		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Volume %d attached to Server %s\n", volume.ID, server.Name)
		return nil
	},
}
