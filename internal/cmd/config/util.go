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

func getPreference(key string) (config.IOption, error) {
	opt, ok := config.Options[key]
	if !ok || !opt.HasFlags(config.OptionFlagPreference) {
		return nil, fmt.Errorf("unknown preference: %s", key)
	}
	return opt, nil
}

func getOptionNames(flags config.OptionFlag) []string {
	var names []string
	for name, opt := range config.Options {
		if opt.HasFlags(flags) {
			names = append(names, name)
		}
	}
	return names
}

func suggestOptionCompletions(_ *cobra.Command, args []string) []string {
	var comps []string
	if opt, ok := config.Options[args[0]]; ok {
		comps = opt.Completions()
	}
	return comps
}
