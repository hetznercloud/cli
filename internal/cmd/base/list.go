package base

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ListCmd allows defining commands for listing resources
type ListCmd struct {
	SortOption         *config.Option[[]string]
	ResourceNamePlural string // e.g. "Servers"
	JSONKeyGetByName   string // e.g. "Servers"
	DefaultColumns     []string
	Fetch              func(state.State, *pflag.FlagSet, hcloud.ListOpts, []string) ([]interface{}, error)
	AdditionalFlags    func(*cobra.Command)
	OutputTable        func(t *output.Table, client hcapi2.Client)
	Schema             func([]interface{}) interface{}
}

// CobraCommand creates a command that can be registered with cobra.
func (lc *ListCmd) CobraCommand(s state.State) *cobra.Command {
	t := output.NewTable(io.Discard)
	lc.OutputTable(t, s.Client())
	outputColumns := t.Columns()

	cmd := &cobra.Command{
		Use:   "list [options]",
		Short: fmt.Sprintf("List %s", lc.ResourceNamePlural),
		Long: util.ListLongDescription(
			fmt.Sprintf("Displays a list of %s.", lc.ResourceNamePlural),
			outputColumns,
		),
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               s.EnsureToken,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return lc.Run(s, cmd)
		},
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(outputColumns), output.OptionJSON(), output.OptionYAML())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	if lc.AdditionalFlags != nil {
		lc.AdditionalFlags(cmd)
	}
	cmd.Flags().StringSliceP("sort", "s", []string{}, "Determine the sorting of the result")
	return cmd
}

// Run executes a list command
func (lc *ListCmd) Run(s state.State, cmd *cobra.Command) error {
	outOpts := output.FlagsForCommand(cmd)

	quiet, err := config.OptionQuiet.Get(s.Config())
	if err != nil {
		return err
	}

	var sorts []string
	if cmd.Flags().Changed("sort") {
		if lc.SortOption == nil {
			_, _ = fmt.Fprintln(os.Stderr, "Warning: resource does not support sorting. Ignoring --sort flag.")
		} else {
			sorts, _ = cmd.Flags().GetStringSlice("sort")
		}
	} else if lc.SortOption != nil {
		var err error
		sorts, err = lc.SortOption.Get(s.Config())
		if err != nil {
			return err
		}
	}

	labelSelector, _ := cmd.Flags().GetString("selector")
	listOpts := hcloud.ListOpts{
		LabelSelector: labelSelector,
		PerPage:       50,
	}

	resources, err := lc.Fetch(s, cmd.Flags(), listOpts, sorts)
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
		schema := lc.Schema(resources)
		if outOpts.IsSet("json") {
			return util.DescribeJSON(out, schema)
		}
		return util.DescribeYAML(out, schema)
	}

	cols := lc.DefaultColumns
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	t := output.NewTable(out)
	lc.OutputTable(t, s.Client())
	if !outOpts.IsSet("noheader") {
		t.WriteHeader(cols)
	}
	for _, resource := range resources {
		t.Write(cols, resource)
	}
	return t.Flush()
}
