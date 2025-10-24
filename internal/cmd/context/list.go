package context

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

type presentation struct {
	Name   string
	Token  string
	Active string
}

func NewListCommand(s state.State) *cobra.Command {
	cols := newListOutputTable(io.Discard).Columns()
	cmd := &cobra.Command{
		Use:   "list [options]",
		Short: "List contexts",
		Long: util.ListLongDescription(
			"Displays a list of contexts.",
			cols,
		),
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(cols))
	return cmd
}

func runList(s state.State, cmd *cobra.Command, _ []string) error {
	outOpts := output.FlagsForCommand(cmd)

	cols := []string{"active", "name"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := newListOutputTable(cmd.OutOrStdout())
	if err := tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	cfg := s.Config()
	for _, context := range cfg.Contexts() {
		presentation := presentation{
			Name:   context.Name(),
			Token:  context.Token(),
			Active: " ",
		}
		if context == cfg.ActiveContext() {
			presentation.Active = "*"
		}

		tw.Write(cols, presentation)
	}
	return tw.Flush()
}

func newListOutputTable(w io.Writer) *output.Table[presentation] {
	return output.NewTable[presentation](w).
		AddAllowedFields(presentation{}).
		RemoveAllowedField("token")
}
