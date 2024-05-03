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
		ok    bool
		err   error
		ctx   config.Context
		prefs config.Preferences
	)

	if global {
		prefs = s.Config().Preferences()
	} else {
		ctx = s.Config().ActiveContext()
		if ctx == nil {
			return fmt.Errorf("no active context (use --global flag to unset a global option)")
		}
		prefs = ctx.Preferences()
	}

	key := args[0]
	if ok, err = prefs.Unset(key); err != nil {
		return err
	}

	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: key '%s' was not set\n", key)
	}
	if ctx == nil {
		cmd.Printf("Unset '%s' globally\n", key)
	} else {
		cmd.Printf("Unset '%s' in context '%s'\n", key, ctx.Name())
	}
	return s.Config().Write(nil)
}
