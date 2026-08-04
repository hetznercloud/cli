package hcapi2

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ImageClient embeds the Hetzner Cloud Image client and provides some
// additional helper functions.
type ImageClient interface {
	hcloud.IImageClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
}

func NewImageClient(client hcloud.IImageClient) ImageClient {
	return &imageClient{
		IImageClient: client,
	}
}

type imageClient struct {
	hcloud.IImageClient
}

// Names obtains a list of available images. It returns nil if image names
// could not be fetched.
func (c *imageClient) Names(ctx context.Context) ([]string, error) {
	images, err := c.AllWithOpts(ctx, hcloud.ImageListOpts{IncludeDeprecated: true})
	if err != nil {
		return nil, err
	}
	return resourceNames(images, func(image *hcloud.Image) int64 { return image.ID }, func(image *hcloud.Image) string { return image.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels assigned to
// the Image with the passed id.
func (c *imageClient) LabelKeys(ctx context.Context, id string) ([]string, error) {
	imgID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid image ID %q: %w", id, err)
	}
	image, _, err := c.GetByID(ctx, imgID)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, nil
	}
	return labelKeys(image.Labels), nil
}
