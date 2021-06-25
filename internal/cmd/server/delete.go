package server

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"

	"github.com/spf13/cobra"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Server",
	ShortDescription:     "Delete a server",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Server().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Server().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		server := resource.(*hcloud.Server)
		if _, err := client.Server().Delete(ctx, server); err != nil {
			return err
		}
		return nil
	},
}
