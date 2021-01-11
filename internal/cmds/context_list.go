package cmds

import (
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

var contextListTableOutput *tableOutput

type ContextPresentation struct {
	Name   string
	Token  string
	Active string
}

func init() {
	contextListTableOutput = newTableOutput().
		AddAllowedFields(ContextPresentation{}).
		RemoveAllowedField("token")
}

func newContextListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List contexts",
		Long: util.ListLongDescription(
			"Displays a list of contexts.",
			contextListTableOutput.Columns(),
		),
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runContextList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(contextListTableOutput.Columns()))
	return cmd
}

func runContextList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	cols := []string{"active", "name"}
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
			Name:   context.Name,
			Token:  context.Token,
			Active: " ",
		}
		if cli.Config.ActiveContext != nil && cli.Config.ActiveContext.Name == context.Name {
			presentation.Active = "*"
		}

		tw.Write(cols, presentation)
	}
	tw.Flush()
	return nil
}
