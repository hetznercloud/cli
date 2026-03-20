package context

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

type presentation struct {
	Name   string
	Token  string
	Active string
}

type schemaPresentation struct {
	Name        string         `json:"name"`
	Token       string         `json:"token,omitempty"`
	Active      bool           `json:"active"`
	Preferences map[string]any `json:"preferences,omitempty"`
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

	cmd.Flags().Bool("allow-sensitive", false, "Allow showing sensitive values in JSON/YAML output (true, false)")

	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(cols), output.OptionJSON(), output.OptionYAML())
	return cmd
}

func runList(s state.State, cmd *cobra.Command, _ []string) error {
	allowSensitive, _ := cmd.Flags().GetBool("allow-sensitive")

	cfg := s.Config()
	outOpts := output.FlagsForCommand(cmd)

	cols := []string{"active", "name"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	quiet, err := config.OptionQuiet.Get(s.Config())
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if quiet {
		// If we are in quiet mode, we saved the original output in cmd.errWriter. We can now restore it.
		out = cmd.ErrOrStderr()
	}

	isSchema := outOpts.IsSet("json") || outOpts.IsSet("yaml")
	if isSchema {
		contexts, activeCtx := cfg.Contexts(), cfg.ActiveContext()
		schema := make([]schemaPresentation, 0, len(contexts))
		for _, ctx := range contexts {
			pres := schemaPresentation{
				Name:        ctx.Name(),
				Active:      ctx == activeCtx,
				Preferences: ctx.Preferences(),
			}
			if allowSensitive {
				pres.Token = ctx.Token()
			}
			schema = append(schema, pres)
		}
		if outOpts.IsSet("json") {
			return util.DescribeJSON(out, schema)
		}
		return util.DescribeYAML(out, schema)
	}

	tw := newListOutputTable(out)
	warnings, err := tw.ValidateColumns(cols)
	if err != nil {
		return err
	}
	for _, warning := range warnings {
		cmd.PrintErrln("Warning:", warning)
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
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
