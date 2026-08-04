package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewDeleteCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete <context>",
		Short:                 "Delete a context",
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidates(config.ContextNames(s.Config())...)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runDelete),
	}
	return cmd
}

func runDelete(s state.State, cmd *cobra.Command, args []string) error {
	name := args[0]
	cfg := s.Config()
	context := config.ContextByName(cfg, name)
	if context == nil {
		return fmt.Errorf("context not found: %v", name)
	}
	if cfg.ActiveContext() == context {
		cmd.PrintErrln("Warning: You deleted the currently active context. Please select a new active context.")
		if err := cfg.SetActiveContext(nil); err != nil {
			return err
		}
	}
	if err := config.RemoveContext(cfg, context); err != nil {
		return err
	}
	return cfg.Write(nil)
}
