package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
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
	Fetch            func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	PrintText        func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error
}

// CobraCommand creates a command that can be registered with cobra.
func (dc *DescribeCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("describe [FLAGS] %s", strings.ToUpper(dc.ResourceNameSingular)),
		Short:                 dc.ShortDescription,
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(dc.NameSuggestions(client))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(tokenEnsurer.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dc.Run(ctx, client, cmd, args)
		},
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	if dc.AdditionalFlags != nil {
		dc.AdditionalFlags(cmd)
	}
	return cmd
}

// Run executes a describe command.
func (dc *DescribeCmd) Run(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	idOrName := args[0]
	resource, resp, err := dc.Fetch(ctx, client, cmd, idOrName)
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
		return dc.describeJSON(resp.Body)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(resource, outputFlags["format"][0])
	default:
		return dc.PrintText(ctx, client, cmd, resource)
	}
}

func (dc *DescribeCmd) describeJSON(body io.ReadCloser) error {
	var data map[string]interface{}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return err
	}
	if resource, ok := data[dc.JSONKeyGetByID]; ok {
		return util.DescribeJSON(resource)
	}
	if resources, ok := data[dc.JSONKeyGetByName].([]interface{}); ok {
		// We check whether we got a resource at all above (see reflect-based nil check), so it's
		// ok to assume there's an element in resources.
		return util.DescribeJSON(resources[0])
	}
	return fmt.Errorf("got invalid JSON response")
}
