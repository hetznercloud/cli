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
			cmpl.SuggestCandidates(getOptionNames(config.OptionFlagPreference|config.OptionFlagSlice)...),
			cmpl.SuggestCandidatesCtx(suggestOptionCompletions),
		)),
	}
	cmd.Flags().Bool("global", false, "Set the value globally (for all contexts)")
	return cmd
}

func runAdd(s state.State, cmd *cobra.Command, args []string) error {
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

	var added []any
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

	if global {
		cmd.Printf("Added '%v' to '%s' globally\n", added, key)
	} else {
		cmd.Printf("Added '%v' to '%s' in context '%s'\n", added, key, ctx.Name())
	}
	return s.Config().Write(nil)
}
