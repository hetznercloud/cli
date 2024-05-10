package config

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewSetCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set <key> <value>...",
		Short:                 "Set a configuration value",
		Long:                  "Set a configuration value. For a list of all available configuration options, run 'hcloud help config'.",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runSet),
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

func runSet(s state.State, cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")

	var (
		val   any
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

	switch t := opt.T().(type) {
	case bool:
		if len(values) != 1 {
			return fmt.Errorf("expected exactly one value")
		}
		value := values[0]
		switch strings.ToLower(value) {
		case "true", "t", "yes", "y", "1":
			val = true
		case "false", "f", "no", "n", "0":
			val = false
		default:
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case string:
		if len(values) != 1 {
			return fmt.Errorf("expected exactly one value")
		}
		val = values[0]
	case time.Duration:
		if len(values) != 1 {
			return fmt.Errorf("expected exactly one value")
		}
		value := values[0]
		var err error
		val, err = time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration value: %s", value)
		}
	case []string:
		newVal := values[:]
		slices.Sort(newVal)
		newVal = slices.Compact(newVal)
		val = newVal
	default:
		return fmt.Errorf("unsupported type %T", t)
	}

	prefs.Set(key, val)

	if util.IsNil(ctx) {
		cmd.Printf("Set '%s' to '%v' globally\n", key, val)
	} else {
		cmd.Printf("Set '%s' to '%v' in context '%s'\n", key, val, ctx.Name())
	}
	return s.Config().Write(nil)
}
