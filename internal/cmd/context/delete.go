package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
)

func newDeleteCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] NAME",
		Short:                 "Delete a context",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(s.Config().ContextNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runDelete),
	}
	return cmd
}

func runDelete(s state.State, _ *cobra.Command, args []string) error {
	name := args[0]
	cfg := s.Config()
	context := cfg.ContextByName(name)
	if context == nil {
		return fmt.Errorf("context not found: %v", name)
	}
	if cfg.ActiveContext() == context {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: You are deleting the currently active context. Please select a new active context.")
		cfg.SetActiveContext(nil)
	}
	cfg.RemoveContext(context)
	return cfg.Write()
}
