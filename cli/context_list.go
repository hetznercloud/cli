package cli

import (
	"github.com/spf13/cobra"
)

var contextListTableOutput *tableOutput

type ContextPresentation struct {
	Name    string
	Token   string
	Current string
}

func init() {
	contextListTableOutput = newTableOutput().
		AddAllowedFields(ContextPresentation{}).
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

	cols := []string{"current", "name"}
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
		presentation := ContextPresentation{
			Name:    context.Name,
			Token:   context.Token,
			Current: " ",
		}
		if cli.Config.ActiveContext.Name == context.Name {
			presentation.Current = "*"
		}

		tw.Write(cols, presentation)
	}
	tw.Flush()
	return nil
}
