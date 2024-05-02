package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewAddCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add <key> <value>...",
		Short:                 "Set a configuration value",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runAdd),
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
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				var comps []string
				if opt, ok := config.Options[args[0]]; ok {
					comps = opt.Completions()
				}
				return comps
			}),
		)),
	}
	cmd.Flags().Bool("global", false, "Set the value globally (for all contexts)")
	return cmd
}

func runAdd(s state.State, cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")

	var prefs config.Preferences

	if global {
		prefs = s.Config().Preferences()
	} else {
		ctx := s.Config().ActiveContext()
		if reflect.ValueOf(ctx).IsNil() {
			return fmt.Errorf("no active context (use --global flag to set a global option)")
		}
		prefs = ctx.Preferences()
	}

	key, values := args[0], args[1:]
	if err := prefs.Add(key, values); err != nil {
		return err
	}

	return s.Config().Write(os.Stdout)
}
