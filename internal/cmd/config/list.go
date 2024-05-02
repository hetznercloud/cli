package config

import (
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

var outputColumns = []string{"key", "value", "origin"}

func NewListCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List configuration values",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runList),
	}
	cmd.Flags().BoolP("all", "a", false, "Also show default values")
	cmd.Flags().BoolP("global", "g", false, "Only show global values")
	cmd.Flags().Bool("allow-sensitive", false, "Allow showing sensitive values")

	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(outputColumns), output.OptionJSON(), output.OptionYAML())
	return cmd
}

func runList(s state.State, cmd *cobra.Command, _ []string) error {
	all, _ := cmd.Flags().GetBool("all")
	global, _ := cmd.Flags().GetBool("global")
	allowSensitive, _ := cmd.Flags().GetBool("allow-sensitive")
	outOpts := output.FlagsForCommand(cmd)

	if global {
		// set context to nil and then reload
		config.OptionContext.OverrideAny(s.Config(), nil)
		s.Config().Reset()
		if err := config.ReadConfig(s.Config(), nil); err != nil {
			return err
		}
	}

	type option struct {
		Key    string `json:"key"`
		Value  any    `json:"value"`
		Origin string `json:"origin"`
	}

	var options []option
	for name, opt := range config.Options {
		val := opt.GetAsAny(s.Config())
		if opt.HasFlag(config.OptionFlagSensitive) && !allowSensitive {
			val = "[redacted]"
		}
		if !all && !opt.Changed(s.Config()) {
			continue
		}
		options = append(options, option{name, val, originToString(s.Config().Viper().Origin(name))})
	}

	// Sort options for reproducible output
	slices.SortFunc(options, func(a, b option) int {
		return strings.Compare(a.Key, b.Key)
	})

	if outOpts.IsSet("json") || outOpts.IsSet("yaml") {
		schema := util.Wrap("options", options)
		if outOpts.IsSet("json") {
			return util.DescribeJSON(schema)
		} else {
			return util.DescribeYAML(schema)
		}
	}

	cols := outputColumns
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	t := output.NewTable()
	t.AddAllowedFields(option{})
	if !outOpts.IsSet("noheader") {
		t.WriteHeader(cols)
	}
	for _, opt := range options {
		t.Write(cols, opt)
	}
	return t.Flush()
}

func originToString(orig viper.ValueOrigin) string {
	switch orig {
	case viper.ValueOriginFlag:
		return "flag"
	case viper.ValueOriginEnv:
		return "environment"
	case viper.ValueOriginConfig:
		return "config file"
	case viper.ValueOriginKVStore:
		return "key-value store"
	case viper.ValueOriginOverride:
		return "override"
	default:
		return "default"
	}
}
