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
		ValidArgsFunction:     cmpl.NoFileCompletion(cmpl.SuggestCandidates(getOptionNames(config.OptionFlagPreference)...)),
	}
	cmd.Flags().Bool("global", false, "Unset the value globally (for all contexts) (true, false)")
	return cmd
}

func runUnset(s state.State, cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")

	ctx, prefs, err := getPreferences(s.Config(), global)
	if err != nil {
		return err
	}

	key := args[0]
	if _, err = getPreference(key); err != nil {
		return err
	}

	if !prefs.Unset(key) {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: key '%s' was not set\n", key)
	}
	if global {
		cmd.Printf("Unset '%s' globally\n", key)
	} else {
		cmd.Printf("Unset '%s' in context '%s'\n", key, ctx.Name())
	}
	return s.Config().Write(nil)
}
