package primaryip

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Primary IP",
	ShortDescription:     "Delete a Primary IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.PrimaryIP().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, _ state.ActionWaiter, cmd *cobra.Command, resource interface{}) error {
		primaryIP := resource.(*hcloud.PrimaryIP)
		if _, err := client.PrimaryIP().Delete(ctx, primaryIP); err != nil {
			return err
		}
		return nil
	},
}
