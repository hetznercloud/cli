package volume

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

var AttachCommand = base.Cmd{
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
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		volume, _, err := client.Volume().Get(ctx, args[0])
		if err != nil {
			return err
		}
		if volume == nil {
			return fmt.Errorf("volume not found: %s", args[0])
		}

		serverIDOrName, _ := cmd.Flags().GetString("server")
		server, _, err := client.Server().Get(ctx, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIDOrName)
		}
		automount, _ := cmd.Flags().GetBool("automount")
		action, _, err := client.Volume().AttachWithOpts(ctx, volume, hcloud.VolumeAttachOpts{
			Server:    server,
			Automount: &automount,
		})

		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		fmt.Printf("Volume %d attached to server %s\n", volume.ID, server.Name)
		return nil
	},
}
