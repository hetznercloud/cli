package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewRenameCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rename <context> <name>",
		Short:                 "Rename a context",
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidates(config.ContextNames(s.Config())...)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runRename),
	}
	return cmd
}

func runRename(s state.State, _ *cobra.Command, args []string) error {
	originalName, newName := args[0], args[1]
	cfg := s.Config()
	context := config.ContextByName(cfg, originalName)
	if context == nil {
		return fmt.Errorf("context not found: %v", originalName)
	}
	isActive := cfg.ActiveContext() == context
	if config.ContextByName(cfg, newName) != nil {
		return fmt.Errorf("context with name %v already exists", newName)
	}
	config.RenameContext(context, newName)
	if isActive {
		// re-set the active context to ensure the name is updated
		cfg.SetActiveContext(context)
	}
	return cfg.Write(nil)
}
