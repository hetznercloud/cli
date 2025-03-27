package context

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewUnsetCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "unset",
		Short:                 "Unset used context",
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestNothing(),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runUnset),
	}
	return cmd
}

func runUnset(s state.State, _ *cobra.Command, _ []string) error {
	cfg := s.Config()
	cfg.SetActiveContext(nil)
	return cfg.Write(nil)
}
