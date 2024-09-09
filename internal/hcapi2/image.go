package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ImageClient embeds the Hetzner Cloud Image client and provides some
// additional helper functions.
type ImageClient interface {
	hcloud.IImageClient
	Names() []string
	LabelKeys(string) []string
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
func (c *imageClient) Names() []string {
	imgs, err := c.AllWithOpts(context.Background(), hcloud.ImageListOpts{IncludeDeprecated: true})
	if err != nil || len(imgs) == 0 {
		return nil
	}
	names := make([]string, len(imgs))
	for i, img := range imgs {
		name := img.Name
		if name == "" {
			name = strconv.FormatInt(img.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels assigned to
// the Image with the passed id.
func (c *imageClient) LabelKeys(id string) []string {
	imgID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil
	}
	img, _, err := c.GetByID(context.Background(), imgID)
	if err != nil || img == nil || len(img.Labels) == 0 {
		return nil
	}
	return labelKeys(img.Labels)
}
