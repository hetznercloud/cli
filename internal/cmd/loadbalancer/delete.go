package loadbalancer

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Delete a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.LoadBalancer().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, _ state.ActionWaiter, cmd *cobra.Command, resource interface{}) error {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		if _, err := client.LoadBalancer().Delete(ctx, loadBalancer); err != nil {
			return err
		}
		return nil
	},
}
