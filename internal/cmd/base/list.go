package base

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Listable is an interface that defines the methods required for a resource to be listed.
// It is needed because ListCmd is a generic type, and we don't always know the concrete type of the resource.
// See [all.ListCmd]
type Listable interface {
	GetResourceNamePlural() string
	GetJSONKeyGetByName() string
	GetDefaultColumns() []string
	NewOutputTable(io.Writer, hcapi2.Client) *output.Table[any]
	FetchAny(state.State, *pflag.FlagSet, hcloud.ListOpts, []string) ([]any, error)
	SchemaAny(any) any
}

// ListCmd allows defining commands for listing resources
// T is the type of the resource that is listed, e.g. *hcloud.Server
// S is the type of the schema that is returned, e.g. schema.Server
type ListCmd[T any, S any] struct {
	SortOption         *config.Option[[]string]
	ResourceNamePlural string // e.g. "Servers"
	JSONKeyGetByName   string // e.g. "Servers"
	DefaultColumns     []string
	Fetch              func(state.State, *pflag.FlagSet, hcloud.ListOpts, []string) ([]T, error)
	// Can be set in case the resource has more than a single identifier that is used in the positional arguments.
	// See [ListCmd.PositionalArgumentOverride].
	FetchWithArgs   func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string, args []string) ([]T, error)
	AdditionalFlags func(*cobra.Command)
	OutputTable     func(t *output.Table[T], client hcapi2.Client)
	Schema          func(T) S

	// In case the resource does not have a single identifier that matches [ListCmd.ResourceNamePlural], this field
	// can be set to define the list of positional arguments.
	// For example, passing:
	//     []string{"a", "b", "c"}
	// Would result in the usage string:
	//     <a> <b> <c>
	PositionalArgumentOverride []string

	// Can be set if auto-completion is needed (usually if [ListCmd.FetchWithArgs] is used)
	ValidArgsFunction func(client hcapi2.Client) cobra.CompletionFunc

	// Experimental is a function that will be used to mark the command as experimental.
	Experimental func(state.State, *cobra.Command) *cobra.Command
}

// CobraCommand creates a command that can be registered with cobra.
func (lc *ListCmd[T, S]) CobraCommand(s state.State) *cobra.Command {
	t := output.NewTable[T](io.Discard)
	lc.OutputTable(t, s.Client())
	outputColumns := t.Columns()

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("list [options]%s", listPositionalArguments(lc.PositionalArgumentOverride)),
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
			return lc.Run(s, cmd, args)
		},
	}

	if lc.ValidArgsFunction != nil {
		cmd.ValidArgsFunction = cmpl.SuggestArgs(lc.ValidArgsFunction(s.Client()))
	}

	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(outputColumns), output.OptionJSON(), output.OptionYAML())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	if lc.AdditionalFlags != nil {
		lc.AdditionalFlags(cmd)
	}
	cmd.Flags().StringSliceP("sort", "s", []string{}, "Determine the sorting of the result")
	if lc.Experimental != nil {
		cmd = lc.Experimental(s, cmd)
	}
	return cmd
}

// Run executes a list command
func (lc *ListCmd[T, S]) Run(s state.State, cmd *cobra.Command, args []string) error {
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

	var resources []T
	if lc.FetchWithArgs != nil {
		resources, err = lc.FetchWithArgs(s, cmd.Flags(), listOpts, sorts, args)
	} else {
		resources, err = lc.Fetch(s, cmd.Flags(), listOpts, sorts)
	}
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
		schema := make([]S, 0, len(resources))
		for _, resource := range resources {
			schema = append(schema, lc.Schema(resource))
		}
		if outOpts.IsSet("json") {
			return util.DescribeJSON(out, schema)
		}
		return util.DescribeYAML(out, schema)
	}

	cols := lc.DefaultColumns
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	t := output.NewTable[T](out)
	lc.OutputTable(t, s.Client())

	warnings, _ := t.ValidateColumns(cols)
	// invalid columns are already checked in output.validateOutputFlag(), we only need the warnings here
	for _, warning := range warnings {
		cmd.PrintErrln("Warning:", warning)
	}

	if !outOpts.IsSet("noheader") {
		t.WriteHeader(cols)
	}
	for _, resource := range resources {
		t.Write(cols, resource)
	}
	return t.Flush()
}

func listPositionalArguments(positionalArgumentOverride []string) string {
	if len(positionalArgumentOverride) == 0 {
		return ""
	}

	return " " + positionalArguments("", positionalArgumentOverride)
}

func (lc *ListCmd[T, S]) GetResourceNamePlural() string {
	return lc.ResourceNamePlural
}

func (lc *ListCmd[T, S]) GetJSONKeyGetByName() string {
	return lc.JSONKeyGetByName
}

func (lc *ListCmd[T, S]) GetDefaultColumns() []string {
	return lc.DefaultColumns
}

func (lc *ListCmd[T, S]) FetchAny(s state.State, fs *pflag.FlagSet, opts hcloud.ListOpts, sorts []string) ([]any, error) {
	resources, err := lc.Fetch(s, fs, opts, sorts)
	return util.ToAnySlice(resources), err
}

func (lc *ListCmd[T, S]) SchemaAny(resource any) any {
	if res, ok := resource.(T); ok {
		return lc.Schema(res)
	}
	return nil
}

func (lc *ListCmd[T, S]) NewOutputTable(out io.Writer, client hcapi2.Client) *output.Table[any] {
	t := output.NewTable[T](out)
	lc.OutputTable(t, client)
	return t.ToAny()
}
