package firewall

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "firewall",
	ShortDescription:     "Delete a firewall",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Firewall().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, resource interface{}) error {
		firewall := resource.(*hcloud.Firewall)
		if _, err := client.Firewall().Delete(ctx, firewall); err != nil {
			return err
		}
		return nil
	},
}
