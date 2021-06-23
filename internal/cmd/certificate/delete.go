package certificate

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "certificate",
	ShortDescription:     "Delete a certificate",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Certificate().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, resource interface{}) error {
		certificate := resource.(*hcloud.Certificate)
		if _, err := client.Certificate().Delete(ctx, certificate); err != nil {
			return err
		}
		return nil
	},
}
