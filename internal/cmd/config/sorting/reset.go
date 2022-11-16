package sorting

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newResetCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "reset COMMAND",
		Short:                 "Reset to the application default sorting order for a command",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runReset),
	}

	return cmd
}

func runReset(cli *state.State, cmd *cobra.Command, args []string) error {
	command := strings.TrimSpace(args[0])
	if command == "" {
		return errors.New("invalid command")
	}

	if cli.Config.SubcommandDefaults == nil {
		return nil
	}

	defaults := cli.Config.SubcommandDefaults[command]
	if defaults != nil {
		defaults.Sorting = nil
	}

	cli.Config.SubcommandDefaults[command] = defaults

	if err := cli.WriteConfig(); err != nil {
		return err
	}

	fmt.Printf("Reset sorting to the default sorting order for command '%s list'\n", command)

	return nil
}
