package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func newUseCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "use <context>",
		Short:                 "Use a context",
		Args:                  util.ValidateExact,
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidates(config.ContextNames(s.Config())...)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runUse),
	}
	return cmd
}

func runUse(s state.State, _ *cobra.Command, args []string) error {
	if os.Getenv("HCLOUD_TOKEN") != "" {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: HCLOUD_TOKEN is set. The active context will have no effect.")
	}
	name := args[0]
	cfg := s.Config()
	context := config.ContextByName(cfg, name)
	if context == nil {
		return fmt.Errorf("context not found: %v", name)
	}
	cfg.SetActiveContext(context)
	return cfg.Write()
}
