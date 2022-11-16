package defaultcolumns

import (
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

var listTableOutput *output.Table

func init() {
	listTableOutput = output.NewTable().
		AddAllowedFields(ColumnsPresentation{})
}

type ColumnsPresentation struct {
	Command string
	Columns string
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [COMMAND]",
		Short:                 "See the default sorting order for a command",
		Args:                  cobra.MaximumNArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runList),
	}

	return cmd
}

func runList(cli *state.State, cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return runListCommand(cli, cmd, args)
	}

	return runListAll(cli, cmd, args)
}

func runListAll(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	cols := []string{"command", "columns"}

	tw := listTableOutput
	if err := tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}

	for command, defaults := range cli.Config.SubcommandDefaults {
		if defaults != nil {
			presentation := ColumnsPresentation{
				Command: command,
				Columns: strings.Join(defaults.DefaultColumns, ", "),
			}

			tw.Write(cols, presentation)
		}
	}

	tw.Flush()
	return nil
}

func runListCommand(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	cols := []string{"command", "columns"}

	command := args[0]

	tw := listTableOutput
	if err := tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}

	defaults := cli.Config.SubcommandDefaults[command]

	if defaults != nil {
		presentation := ColumnsPresentation{
			Command: command,
			Columns: strings.Join(defaults.DefaultColumns, ", "),
		}

		tw.Write(cols, presentation)
	}

	tw.Flush()
	return nil
}
