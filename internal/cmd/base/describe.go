package base

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

// DescribeCmd allows defining commands for describing a resource.
type DescribeCmd[T any] struct {
	ResourceNameSingular string // e.g. "Server"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	AdditionalFlags      func(*cobra.Command)
	// Fetch is called to fetch the resource to describe.
	// The first returned interface is the resource itself as a hcloud struct, the second is the schema for the resource.
	Fetch func(s state.State, cmd *cobra.Command, idOrName string) (T, any, error)
	// Can be set in case the resource has more than a single identifier that is used in the positional arguments.
	// See [DescribeCmd.PositionalArgumentOverride].
	FetchWithArgs func(s state.State, cmd *cobra.Command, args []string) (T, any, error)

	PrintText   func(s state.State, cmd *cobra.Command, resource T) error
	GetIDOrName func(resource T) string

	// In case the resource does not have a single identifier that matches [DescribeCmd.ResourceNameSingular], this field
	// can be set to define the list of positional arguments.
	// For example, passing:
	//     []string{"a", "b", "c"}
	// Would result in the usage string:
	//     <a> <b> <c>
	PositionalArgumentOverride []string

	// Can be set if the default [DescribeCmd.NameSuggestions] is not enough. This is usually the case when
	// [DescribeCmd.FetchWithArgs] and [DescribeCmd.PositionalArgumentOverride] is being used.
	ValidArgsFunction func(client hcapi2.Client) []cobra.CompletionFunc

	// ExperimentalF is a function that will be used to mark the command as experimental.
	ExperimentalF func(state.State, *cobra.Command) *cobra.Command
}

// CobraCommand creates a command that can be registered with cobra.
func (dc *DescribeCmd[T]) CobraCommand(s state.State) *cobra.Command {
	var suggestArgs []cobra.CompletionFunc
	switch {
	case dc.NameSuggestions != nil:
		suggestArgs = append(suggestArgs,
			cmpl.SuggestCandidatesF(dc.NameSuggestions(s.Client())),
		)
	case dc.ValidArgsFunction != nil:
		suggestArgs = append(suggestArgs, dc.ValidArgsFunction(s.Client())...)
	default:
		log.Fatalf("describe command %s is missing ValidArgsFunction or NameSuggestions", dc.ResourceNameSingular)
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("describe [options] %s", positionalArguments(dc.ResourceNameSingular, dc.PositionalArgumentOverride)),
		Short:                 dc.ShortDescription,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(suggestArgs...),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dc.Run(s, cmd, args)
		},
	}

	output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML(), output.OptionFormat())

	if dc.AdditionalFlags != nil {
		dc.AdditionalFlags(cmd)
	}

	if dc.ExperimentalF != nil {
		cmd = dc.ExperimentalF(s, cmd)
	}

	return cmd
}

// Run executes a describe command.
func (dc *DescribeCmd[T]) Run(s state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	quiet, err := config.OptionQuiet.Get(s.Config())
	if err != nil {
		return err
	}

	schemaOut := cmd.OutOrStdout()
	isSchema := outputFlags.IsSet("json") || outputFlags.IsSet("yaml")
	if isSchema {
		if quiet {
			// If we are in quiet mode, we saved the original output in cmd.errWriter. We can now restore it.
			schemaOut = cmd.ErrOrStderr()
		} else {
			// We don't want anything other than the schema in stdout, so we set the default to stderr
			cmd.SetOut(os.Stderr)
		}
	}

	var resource T
	var schema any
	if dc.FetchWithArgs != nil {
		resource, schema, err = dc.FetchWithArgs(s, cmd, args)
	} else {
		resource, schema, err = dc.Fetch(s, cmd, args[0])
	}
	if err != nil {
		return err
	}

	// resource is an interface that always has a type, so the interface is never nil
	// (i.e. == nil) is always false.
	if reflect.ValueOf(resource).IsNil() {
		return fmt.Errorf("%s not found: %s", dc.ResourceNameSingular, strings.Join(args, " "))
	}

	switch {
	case outputFlags.IsSet("json"):
		return util.DescribeJSON(schemaOut, schema)
	case outputFlags.IsSet("yaml"):
		return util.DescribeYAML(schemaOut, schema)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(schemaOut, resource, outputFlags["format"][0])
	default:
		return dc.PrintText(s, cmd, resource)
	}
}
