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
	ResourceNamePlural string // e.g. "servers"
	JSONKeyGetByName   string // e.g. "servers"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.Run(s, cmd)
		},
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(outputColumns), output.OptionJSON(), output.OptionYAML())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	if lc.AdditionalFlags != nil {
		lc.AdditionalFlags(cmd)
	}
	cmd.Flags().StringSliceP("sort", "s", []string{"id:asc"}, "Determine the sorting of the result")
	return cmd
}

// Run executes a list command
func (lc *ListCmd) Run(s state.State, cmd *cobra.Command) error {
	outOpts := output.FlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	listOpts := hcloud.ListOpts{
		LabelSelector: labelSelector,
		PerPage:       50,
	}
	sorts, _ := cmd.Flags().GetStringSlice("sort")

	resources, err := lc.Fetch(s, cmd.Flags(), listOpts, sorts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") || outOpts.IsSet("yaml") {
		schema := lc.Schema(resources)
		if outOpts.IsSet("json") {
			return util.DescribeJSON(cmd.OutOrStdout(), schema)
		} else {
			return util.DescribeYAML(cmd.OutOrStdout(), schema)
		}
	}

	cols := lc.DefaultColumns
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	out := cmd.OutOrStdout()
	if config.OptionQuiet.Get(s.Config()) {
		// if the quiet option is set, write to stdout anyway, since outputting nothing would not make sense
		out = os.Stdout
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
