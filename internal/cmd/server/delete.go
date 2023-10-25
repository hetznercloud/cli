package server

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Server",
	ShortDescription:     "Delete a server",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Server().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Server().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, resource interface{}) error {
		server := resource.(*hcloud.Server)
		result, _, err := client.Server().DeleteWithResult(ctx, server)
		if err != nil {
			return err
		}

		if err := actionWaiter.ActionProgress(ctx, result.Action); err != nil {
			return err
		}

		return nil
	},
}
