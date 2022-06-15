package base

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

// ListCmd allows defining commands for listing resources
type ListCmd struct {
	ResourceNamePlural string // e.g. "servers"
	DefaultColumns     []string
	Fetch              func(context.Context, hcapi2.Client, *cobra.Command, hcloud.ListOpts, []string) ([]interface{}, error)
	AdditionalFlags    func(*cobra.Command)
	OutputTable        func(client hcapi2.Client) *output.Table
	JSONSchema         func([]interface{}) interface{}
}

// CobraCommand creates a command that can be registered with cobra.
func (lc *ListCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer,
) *cobra.Command {
	outputColumns := lc.OutputTable(client).Columns()

	cmd := &cobra.Command{
		Use:   "list [FlAGS]",
		Short: fmt.Sprintf("List %s", lc.ResourceNamePlural),
		Long: util.ListLongDescription(
			fmt.Sprintf("Displays a list of %s.", lc.ResourceNamePlural),
			outputColumns,
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               tokenEnsurer.EnsureToken,
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.Run(ctx, client, cmd)
		},
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(outputColumns), output.OptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	if lc.AdditionalFlags != nil {
		lc.AdditionalFlags(cmd)
	}
	cmd.Flags().StringSliceP("sort", "s", []string{"id:asc"}, "Determine the sorting of the result")
	return cmd
}

// Run executes a list command
func (lc *ListCmd) Run(ctx context.Context, client hcapi2.Client, cmd *cobra.Command) error {
	outOpts := output.FlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	listOpts := hcloud.ListOpts{
		LabelSelector: labelSelector,
		PerPage:       50,
	}
	sorts, _ := cmd.Flags().GetStringSlice("sort")

	resources, err := lc.Fetch(ctx, client, cmd, listOpts, sorts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		return util.DescribeJSON(lc.JSONSchema(resources))
	}

	cols := lc.DefaultColumns
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	table := lc.OutputTable(client)
	if !outOpts.IsSet("noheader") {
		table.WriteHeader(cols)
	}
	for _, resource := range resources {
		table.Write(cols, resource)
	}
	table.Flush()

	return nil
}
