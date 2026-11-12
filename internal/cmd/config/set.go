package config

import (
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
			cmpl.SuggestCandidates(getOptionNames(config.OptionFlagPreference)...),
			cmpl.SuggestCandidatesCtx(suggestOptionCompletions),
		)),
	}
	cmd.Flags().Bool("global", false, "Set the value globally (for all contexts)")
	return cmd
}

func runSet(s state.State, cmd *cobra.Command, args []string) error {
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

	val, err := opt.Parse(values)
	if err != nil {
		return err
	}

	prefs.Set(key, val)

	if global {
		cmd.Printf("Set '%s' to '%v' globally\n", key, val)
	} else {
		cmd.Printf("Set '%s' to '%v' in context '%s'\n", key, val, ctx.Name())
	}
	return s.Config().Write(nil)
}
