package server

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var ChangeTypeCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "change-type [FLAGS] SERVER SERVERTYPE",
			Short: "Change type of a server",
			Args:  cobra.ExactArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Server().Names),
				cmpl.SuggestCandidatesF(client.ServerType().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("keep-disk", false, "Keep disk size of current server type. This enables downgrading the server.")
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

		serverTypeIDOrName := args[1]
		serverType, _, err := client.ServerType().Get(ctx, serverTypeIDOrName)
		if err != nil {
			return err
		}
		if serverType == nil {
			return fmt.Errorf("server type not found: %s", serverTypeIDOrName)
		}

		if serverType.IsDeprecated() {
			fmt.Print(warningDeprecatedServerType(serverType))
		}

		keepDisk, _ := cmd.Flags().GetBool("keep-disk")
		opts := hcloud.ServerChangeTypeOpts{
			ServerType:  serverType,
			UpgradeDisk: !keepDisk,
		}
		action, _, err := client.Server().ChangeType(ctx, server, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		if opts.UpgradeDisk {
			fmt.Printf("Server %d changed to type %s\n", server.ID, serverType.Name)
		} else {
			fmt.Printf("Server %d changed to type %s (disk size was unchanged)\n", server.ID, serverType.Name)
		}
		return nil
	},
}
