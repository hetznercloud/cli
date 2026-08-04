package floatingip

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var AssignCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "assign <floating-ip> <server>",
			Short: "Assign a Floating IP to a server",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.FloatingIP().Names),
				cmpl.SuggestCandidatesF(client.Server().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		floatingIP, _, err := s.Client().FloatingIP().Get(cmd.Context(), idOrName)
		if err != nil {
			return err
		}
		if floatingIP == nil {
			return fmt.Errorf("Floating IP not found: %v", idOrName)
		}

		serverIDOrName := args[1]
		server, _, err := s.Client().Server().Get(cmd.Context(), serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", serverIDOrName)
		}

		action, _, err := s.Client().FloatingIP().Assign(cmd.Context(), floatingIP, server)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, action); err != nil {
			return err
		}

		cmd.Printf("Floating IP %d assigned to Server %d\n", floatingIP.ID, server.ID)
		return nil
	},
}
