package base

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

// DescribeCmd allows defining commands for describing a resource.
type DescribeCmd struct {
	ResourceNameSingular string // e.g. "server"
	ShortDescription     string
	// key in API response JSON to use for extracting object from response body for JSON output.
	JSONKeyGetByID   string // e.g. "server"
	JSONKeyGetByName string // e.g. "servers"
	NameSuggestions  func(client hcapi2.Client) func() []string
	AdditionalFlags  func(*cobra.Command)
	// Fetch is called to fetch the resource to describe.
	// The first returned interface is the resource itself as a hcloud struct, the second is the schema for the resource.
	Fetch     func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error)
	PrintText func(s state.State, cmd *cobra.Command, resource interface{}) error
}

// CobraCommand creates a command that can be registered with cobra.
func (dc *DescribeCmd) CobraCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("describe [options] <%s>", util.ToKebabCase(dc.ResourceNameSingular)),
		Short:                 dc.ShortDescription,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(dc.NameSuggestions(s.Client()))),
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
	return cmd
}

// Run executes a describe command.
func (dc *DescribeCmd) Run(s state.State, cmd *cobra.Command, args []string) error {
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

	idOrName := args[0]
	resource, schema, err := dc.Fetch(s, cmd, idOrName)
	if err != nil {
		return err
	}

	// resource is an interface that always has a type, so the interface is never nil
	// (i.e. == nil) is always false.
	if reflect.ValueOf(resource).IsNil() {
		return fmt.Errorf("%s not found: %s", dc.ResourceNameSingular, idOrName)
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
