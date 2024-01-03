package floatingip

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var UnassignCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:                   "unassign [FLAGS] FLOATINGIP",
			Short:                 "Unassign a Floating IP",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.FloatingIP().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		floatingIP, _, err := s.FloatingIP().Get(s, idOrName)
		if err != nil {
			return err
		}
		if floatingIP == nil {
			return fmt.Errorf("Floating IP not found: %v", idOrName)
		}

		action, _, err := s.FloatingIP().Unassign(s, floatingIP)
		if err != nil {
			return err
		}

		if err := s.ActionProgress(cmd, s, action); err != nil {
			return err
		}

		cmd.Printf("Floating IP %d unassigned\n", floatingIP.ID)
		return nil
	},
}
