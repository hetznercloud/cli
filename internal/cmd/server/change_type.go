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

var ChangeTypeCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "change-type [--keep-disk] <server> <server-type>",
			Short: "Change type of a server",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Server().Names),
				cmpl.SuggestCandidatesF(client.ServerType().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("keep-disk", false, "Keep disk size of current Server Type. This enables downgrading the server. (true, false)")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}

		serverTypeIDOrName := args[1]
		serverType, _, err := s.Client().ServerType().Get(s, serverTypeIDOrName)
		if err != nil {
			return err
		}
		if serverType == nil {
			return fmt.Errorf("Server Type not found: %s", serverTypeIDOrName)
		}

		cmd.Print(deprecatedServerTypeWarning(serverType, server.Location.Name))

		keepDisk, _ := cmd.Flags().GetBool("keep-disk")
		opts := hcloud.ServerChangeTypeOpts{
			ServerType:  serverType,
			UpgradeDisk: !keepDisk,
		}
		action, _, err := s.Client().Server().ChangeType(s, server, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		if opts.UpgradeDisk {
			cmd.Printf("Server %d changed to type %s\n", server.ID, serverType.Name)
		} else {
			cmd.Printf("Server %d changed to type %s (disk size was unchanged)\n", server.ID, serverType.Name)
		}
		return nil
	},
}
