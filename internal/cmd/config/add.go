package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewAddCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add <key> <value>...",
		Short:                 "Add values to a list",
		Long:                  "Add values to a list. For a list of all available configuration options, run 'hcloud help config'.",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
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

	var (
		added []any
		ctx   config.Context
		prefs config.Preferences
	)

	if global {
		prefs = s.Config().Preferences()
	} else {
		ctx = s.Config().ActiveContext()
		if util.IsNil(ctx) {
			return fmt.Errorf("no active context (use --global flag to set a global option)")
		}
		prefs = ctx.Preferences()
	}

	key, values := args[0], args[1:]
	opt, ok := config.Options[key]
	if !ok || !opt.HasFlag(config.OptionFlagPreference) {
		return fmt.Errorf("unknown preference: %s", key)
	}

	val, _ := prefs.Get(key)
	switch opt.T().(type) {
	case []string:
		before := util.AnyToStringSlice(val)
		newVal := append(before, values...)
		slices.Sort(newVal)
		newVal = slices.Compact(newVal)
		val = newVal
		added = util.ToAnySlice(util.SliceDiff[[]string](newVal, before))
	default:
		return fmt.Errorf("%s is not a list", key)
	}

	prefs.Set(key, val)

	if len(added) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: no new values were added")
	} else if len(added) < len(values) {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: some values were already present or duplicate")
	}

	if util.IsNil(ctx) {
		cmd.Printf("Added '%v' to '%s' globally\n", added, key)
	} else {
		cmd.Printf("Added '%v' to '%s' in context '%s'\n", added, key, ctx.Name())
	}
	return s.Config().Write(nil)
}
