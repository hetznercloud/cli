package base

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DeleteCmd allows defining commands for deleting a resource.
type DeleteCmd struct {
	ResourceNameSingular string // e.g. "server"
	ShortDescription     string
	NameSuggestions      func(client hcapi2.Client) func() []string
	AdditionalFlags      func(*cobra.Command)
	Fetch                func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error)
	Delete               func(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, resource interface{}) error
}

// CobraCommand creates a command that can be registered with cobra.
func (dc *DeleteCmd) CobraCommand(
	ctx context.Context, client hcapi2.Client, tokenEnsurer state.TokenEnsurer, actionWaiter state.ActionWaiter,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("delete [FLAGS] %s", strings.ToUpper(dc.ResourceNameSingular)),
		Short:                 dc.ShortDescription,
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(dc.NameSuggestions(client))),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(tokenEnsurer.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dc.Run(ctx, client, actionWaiter, cmd, args)
		},
	}
	if dc.AdditionalFlags != nil {
		dc.AdditionalFlags(cmd)
	}
	return cmd
}

// Run executes a describe command.
func (dc *DeleteCmd) Run(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, args []string) error {

	idOrName := args[0]
	resource, _, err := dc.Fetch(ctx, client, cmd, idOrName)
	if err != nil {
		return err
	}

	// resource is an interface that always has a type, so the interface is never nil
	// (i.e. == nil) is always false.
	if reflect.ValueOf(resource).IsNil() {
		return fmt.Errorf("%s not found: %s", dc.ResourceNameSingular, idOrName)
	}

	if err := dc.Delete(ctx, client, actionWaiter, cmd, resource); err != nil {
		return fmt.Errorf("deleting %s %s failed: %s", dc.ResourceNameSingular, idOrName, err)
	}
	cmd.Printf("%s %v deleted\n", dc.ResourceNameSingular, idOrName)
	return nil
}
