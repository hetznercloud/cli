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

func NewRemoveCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove <key> <value>...",
		Short:                 "Remove values from a list",
		Long:                  "Remove values from a list. For a list of all available configuration options, run 'hcloud help config'.",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runRemove),
		ValidArgsFunction: cmpl.NoFileCompletion(cmpl.SuggestArgs(
			cmpl.SuggestCandidates(getOptionNames(config.OptionFlagPreference|config.OptionFlagSlice)...),
			cmpl.SuggestCandidatesCtx(suggestOptionCompletions),
		)),
	}
	cmd.Flags().Bool("global", false, "Remove the value(s) globally (for all contexts) (true, false)")
	return cmd
}

func runRemove(s state.State, cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")

	ctx, prefs, err := getPreferences(s.Config(), global)
	if err != nil {
		return err
	}

	key, values := args[0], args[1:]
	opt, err := getPreference(key)
	if err != nil {
		return err
	}

	val, _ := prefs.Get(key)

	var removed []any
	switch opt.T().(type) {
	case []string:
		before := util.AnyToStringSlice(val)
		diff := util.SliceDiff[[]string](before, values)
		val = diff
		removed = util.ToAnySlice(util.SliceDiff[[]string](before, diff))
	default:
		return fmt.Errorf("%s is not a list", key)
	}

	if reflect.ValueOf(val).Len() == 0 {
		prefs.Unset(key)
	} else {
		prefs.Set(key, val)
	}

	if len(removed) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: no values were removed")
	} else if len(removed) < len(values) {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: some values were not removed")
	}

	if global {
		cmd.Printf("Removed '%v' from '%s' globally\n", removed, key)
	} else {
		cmd.Printf("Removed '%v' from '%s' in context '%s'\n", removed, key, ctx.Name())
	}
	return s.Config().Write(nil)
}
