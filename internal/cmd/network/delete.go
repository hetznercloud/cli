package network

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"

	"github.com/spf13/cobra"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Network",
	ShortDescription:     "Delete a network",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Network().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Network().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, _ state.ActionWaiter, cmd *cobra.Command, resource interface{}) error {
		network := resource.(*hcloud.Network)
		if _, err := client.Network().Delete(ctx, network); err != nil {
			return err
		}
		return nil
	},
}
