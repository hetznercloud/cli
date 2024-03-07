package server

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
			Use:   "disable-protection <server> <protection-level>...",
			Short: "Disable resource protection for a server",
			Args:  cobra.MinimumNArgs(2),
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Server().Names),
				cmpl.SuggestCandidates("delete", "rebuild"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		opts, err := getChangeProtectionOpts(false, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, server, false, opts)
	},
}
