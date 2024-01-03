package floatingip

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var DisableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "disable-protection [FLAGS] FLOATINGIP PROTECTIONLEVEL [PROTECTIONLEVEL...]",
			Short: "Disable resource protection for a Floating IP",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.FloatingIP().Names),
				cmpl.SuggestCandidates("delete"),
			),
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

		opts, err := getChangeProtectionOpts(false, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, floatingIP, false, opts)
	},
}
