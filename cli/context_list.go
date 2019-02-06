package cli

import (
	"github.com/spf13/cobra"
)

var contextListTableOutput *tableOutput

func init() {
	contextListTableOutput = newTableOutput().
		AddAllowedFields(ConfigContext{}).
		RemoveAllowedField("token")
}

func newContextListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List contexts",
		Long: listLongDescription(
			"Displays a list of contexts.",
			contextListTableOutput.Columns(),
		),
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runContextList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(contextListTableOutput.Columns()))
	return cmd
}

func runContextList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	cols := []string{"name"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := contextListTableOutput
	if err := tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, context := range cli.Config.Contexts {
		if context.Token == cli.Token {
			context.Name = context.Name + " |ACTIVE|"
		}
		tw.Write(cols, context)
	}
	tw.Flush()
	return nil
}
