package base

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

// Cmd allows defining commands for generic resource-based commands
type Cmd struct {
	BaseCobraCommand func(hcapi2.Client) *cobra.Command
	Run              func(state.State, *cobra.Command, []string) error
}

// CobraCommand creates a command that can be registered with cobra.
func (gc *Cmd) CobraCommand(s state.State) *cobra.Command {
	cmd := gc.BaseCobraCommand(s.Client())

	if cmd.Args == nil {
		cmd.Args = cobra.NoArgs
	}

	cmd.TraverseChildren = true
	cmd.DisableFlagsInUseLine = true

	if cmd.PreRunE != nil {
		cmd.PreRunE = util.ChainRunE(cmd.PreRunE, s.EnsureToken)
	} else {
		cmd.PreRunE = s.EnsureToken
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return gc.Run(s, cmd, args)
	}

	return cmd
}
