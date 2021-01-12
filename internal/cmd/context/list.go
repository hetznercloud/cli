package context

import (
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

var listTableOutput *output.Table

type ContextPresentation struct {
	Name   string
	Token  string
	Active string
}

func init() {
	listTableOutput = output.NewTable().
		AddAllowedFields(ContextPresentation{}).
		RemoveAllowedField("token")
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List contexts",
		Long: util.ListLongDescription(
			"Displays a list of contexts.",
			listTableOutput.Columns(),
		),
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(listTableOutput.Columns()))
	return cmd
}

func runList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	cols := []string{"active", "name"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := listTableOutput
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
