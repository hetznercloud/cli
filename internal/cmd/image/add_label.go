package image

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var addLabelCmd = base.AddLabelCmd{
	ResourceNameSingular: "image",
	ShortDescription:     "Add a label to an image",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Image().Names },
	FetchLabels: func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int, error) {
		image, _, err := client.Image().Get(ctx, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if image == nil {
			return nil, 0, fmt.Errorf("image not found: %s", idOrName)
		}
		return image.Labels, image.ID, nil
	},
	SetLabels: func(ctx context.Context, client hcapi2.Client, id int, labels map[string]string) error {
		opts := hcloud.ImageUpdateOpts{
			Labels: labels,
		}
		_, _, err := client.Image().Update(ctx, &hcloud.Image{ID: id}, opts)
		return err
	},
}
