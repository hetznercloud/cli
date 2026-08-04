package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state/config"
)

func getPreferences(cfg config.Config, global bool) (ctx config.Context, prefs config.Preferences, _ error) {
	if global {
		prefs = cfg.Preferences()
	} else {
		ctx = cfg.ActiveContext()
		if util.IsNil(ctx) {
			return nil, nil, fmt.Errorf("no active context (use --global flag to set a global option)")
		}
		prefs = ctx.Preferences()
	}
	return
}

func getPreference(cfg config.Config, key string) (config.IOption, error) {
	opt, ok := cfg.LookupOption(key)
	if !ok || !opt.HasFlags(config.OptionFlagPreference) {
		return nil, fmt.Errorf("unknown preference: %s", key)
	}
	return opt, nil
}

func getOptionNames(cfg config.Config, flags config.OptionFlag) []string {
	var names []string
	for _, opt := range cfg.Options() {
		if opt.HasFlags(flags) {
			names = append(names, opt.GetName())
		}
	}
	return names
}

func suggestOptionCompletions(cfg config.Config) func(*cobra.Command, []string) []string {
	return func(_ *cobra.Command, args []string) []string {
		if len(args) == 0 {
			return nil
		}
		if opt, ok := cfg.LookupOption(args[0]); ok {
			return opt.Completions()
		}
		return nil
	}
}
