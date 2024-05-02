package config

import (
	"fmt"
	"reflect"

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
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
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

	var prefs config.Preferences

	if global {
		prefs = s.Config().Preferences()
	} else {
		ctx := s.Config().ActiveContext()
		if reflect.ValueOf(ctx).IsNil() {
			return fmt.Errorf("no active context (use --global flag to unset a global option)")
		}
		prefs = ctx.Preferences()
	}

	key := args[0]
	if err := prefs.Unset(key); err != nil {
		return err
	}

	return s.Config().Write(nil)
}
