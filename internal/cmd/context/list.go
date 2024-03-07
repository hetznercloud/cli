package context

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
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

func newListCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [options]",
		Short: "List contexts",
		Long: util.ListLongDescription(
			"Displays a list of contexts.",
			listTableOutput.Columns(),
		),
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(listTableOutput.Columns()))
	return cmd
}

func runList(s state.State, cmd *cobra.Command, _ []string) error {
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
	cfg := s.Config()
	for _, context := range cfg.Contexts() {
		presentation := ContextPresentation{
			Name:   context.Name,
			Token:  context.Token,
			Active: " ",
		}
		if ctx := cfg.ActiveContext(); ctx != nil && ctx.Name == context.Name {
			presentation.Active = "*"
		}

		tw.Write(cols, presentation)
	}
	tw.Flush()
	return nil
}
