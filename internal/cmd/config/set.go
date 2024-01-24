package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func newSetCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set <key> <value>",
		Short:                 "Set a configuration value",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runSet),
	}
	cmd.Flags().Bool("global", false, "Set the value globally (for all contexts)")
	return cmd
}

func runSet(s state.State, cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")

	var prefs config.Preferences

	if global {
		prefs = s.Config().Preferences()
	} else {
		ctx := s.Config().ActiveContext()
		if ctx == nil {
			if ctxName := config.OptionContext.Value(); ctxName != "" {
				return fmt.Errorf("active context \"%s\" not found", ctxName)
			} else {
				return fmt.Errorf("no active context (use --global flag to set a global option)")
			}
		}
		prefs = ctx.Preferences()
	}

	key, value := args[0], args[1]
	if err := prefs.Set(key, value); err != nil {
		return err
	}

	return s.Config().Write(nil)
}
