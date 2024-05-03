package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewGetCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get <key>",
		Short:                 "Get a configuration value",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runGet),
	}
	cmd.Flags().Bool("global", false, "Get the value globally")
	cmd.Flags().Bool("allow-sensitive", false, "Allow showing sensitive values")
	return cmd
}

func runGet(s state.State, cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")
	allowSensitive, _ := cmd.Flags().GetBool("allow-sensitive")

	if global {
		// set context to nil and then reload
		config.OptionContext.OverrideAny(s.Config(), nil)
		s.Config().Reset()
		if err := config.ReadConfig(s.Config(), nil); err != nil {
			return err
		}
	}

	key := args[0]
	var opt config.IOption
	for name, o := range config.Options {
		if name == key {
			opt = o
			break
		}
	}
	if opt == nil {
		return fmt.Errorf("unknown key: %s", key)
	}

	val := opt.GetAsAny(s.Config())
	if opt.HasFlag(config.OptionFlagSensitive) && !allowSensitive {
		return fmt.Errorf("'%s' is sensitive. use --allow-sensitive to show the value", key)
	}
	cmd.Println(val)
	return nil
}
