package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewUnsetCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "unset <key>",
		Short:                 "Unset a configuration value",
		Long:                  "Unset a configuration value. For a list of all available configuration options, run 'hcloud help config'.",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runUnset),
		ValidArgsFunction: cmpl.NoFileCompletion(cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(func() []string {
				var keys []string
				for key, opt := range config.Options {
					if opt.HasFlag(config.OptionFlagPreference) {
						keys = append(keys, key)
					}
				}
				return keys
			}),
		)),
	}
	cmd.Flags().Bool("global", false, "Unset the value globally (for all contexts)")
	return cmd
}

func runUnset(s state.State, cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")

	var (
		ctx   config.Context
		prefs config.Preferences
	)

	if global {
		prefs = s.Config().Preferences()
	} else {
		ctx = s.Config().ActiveContext()
		if util.IsNil(ctx) {
			return fmt.Errorf("no active context (use --global flag to unset a global option)")
		}
		prefs = ctx.Preferences()
	}

	key := args[0]
	opt, ok := config.Options[key]
	if !ok || !opt.HasFlag(config.OptionFlagPreference) {
		return fmt.Errorf("unknown preference: %s", key)
	}

	ok = prefs.Unset(key)

	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: key '%s' was not set\n", key)
	}
	if util.IsNil(ctx) {
		cmd.Printf("Unset '%s' globally\n", key)
	} else {
		cmd.Printf("Unset '%s' in context '%s'\n", key, ctx.Name())
	}
	return s.Config().Write(nil)
}
