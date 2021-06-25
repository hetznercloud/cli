package volume

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var deleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Volume",
	ShortDescription:     "Delete a Volume",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Volume().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Volume().Get(ctx, idOrName)
	},
	Delete: func(ctx context.Context, client hcapi2.Client, resource interface{}) error {
		volume := resource.(*hcloud.Volume)
		if _, err := client.Volume().Delete(ctx, volume); err != nil {
			return err
		}
		return nil
	},
}
