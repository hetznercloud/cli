package defaultcolumns

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newSetCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set COMMAND COLUMNS...",
		Short:                 "Set which columns to show by default for a command",
		Args:                  cobra.MinimumNArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runSet),
	}

	return cmd
}

func runSet(cli *state.State, cmd *cobra.Command, args []string) error {
	command := strings.TrimSpace(args[0])
	if command == "" {
		return errors.New("invalid command")
	}

	if len(args[1:]) == 0 {
		return errors.New("invalid columns")
	}

	columns := make([]string, len(args[1:]))
	for index, columnName := range args[1:] {
		columns[index] = strings.TrimSpace(columnName)
	}

	if cli.Config.SubcommandDefaults == nil {
		cli.Config.SubcommandDefaults = make(map[string]*state.SubcommandDefaults)
	}

	defaults := cli.Config.SubcommandDefaults[command]
	if defaults == nil {
		defaults = &state.SubcommandDefaults{
			DefaultColumns: columns,
		}
	} else {
		defaults.DefaultColumns = columns
	}

	cli.Config.SubcommandDefaults[command] = defaults

	if err := cli.WriteConfig(); err != nil {
		return err
	}

	fmt.Printf("Displaying columns '%s' for command '%s list'\n", strings.Join(columns, ", "), command)

	return nil
}
