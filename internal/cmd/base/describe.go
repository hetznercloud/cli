package base

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
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
		Use:                   fmt.Sprintf("describe [FLAGS] %s", strings.ToUpper(dc.ResourceNameSingular)),
		Short:                 dc.ShortDescription,
		Args:                  cobra.ExactArgs(1),
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

	quiet, _ := cmd.Flags().GetBool("quiet")

	isSchema := outputFlags.IsSet("json") || outputFlags.IsSet("yaml")
	if isSchema && !quiet {
		cmd.SetOut(os.Stderr)
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
		return util.DescribeJSON(schema)
	case outputFlags.IsSet("yaml"):
		return util.DescribeYAML(schema)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(resource, outputFlags["format"][0])
	default:
		return dc.PrintText(s, cmd, resource)
	}
}
