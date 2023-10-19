package floatingip

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Floating IP",
	ShortDescription:     "Delete a Floating IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.FloatingIP().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, _ state.ActionWaiter, cmd *cobra.Command, resource interface{}) error {
		floatingIP := resource.(*hcloud.FloatingIP)
		if _, err := client.FloatingIP().Delete(ctx, floatingIP); err != nil {
			return err
		}
		return nil
	},
}
