package sorting

import (
	"errors"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

var listTableOutput *output.Table

func init() {
	listTableOutput = output.NewTable().
		AddAllowedFields(SortPresentation{})
}

type SortPresentation struct {
	Command string
	Column  string
	Order   string
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

	cols := []string{"command", "column"}

	tw := listTableOutput
	if err := tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}

	for command, defaults := range cli.Config.SubcommandDefaults {
		if defaults != nil {
			presentation := SortPresentation{
				Command: command,
				Column:  strings.Join(defaults.Sorting, ", "),
				Order:   "",
			}

			tw.Write(cols, presentation)
		}
	}

	tw.Flush()
	return nil
}

func runListCommand(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	cols := []string{"column", "order"}

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
		for _, column := range defaults.Sorting {
			order := "asc"

			// handle special case where colum-name includes the : to specify the order
			if strings.Contains(column, ":") {
				columnWithOrdering := strings.Split(column, ":")
				if len(columnWithOrdering) != 2 {
					return errors.New("Column sort syntax invalid")
				}

				column = columnWithOrdering[0]
				order = columnWithOrdering[1]
			}

			presentation := SortPresentation{
				Command: command,
				Column:  column,
				Order:   order,
			}

			tw.Write(cols, presentation)
		}
	}

	tw.Flush()
	return nil
}
