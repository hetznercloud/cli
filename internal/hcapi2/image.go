package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// ImageClient embeds the Hetzner Cloud Image client and provides some
// additional helper functions.
type ImageClient interface {
	ImageClientBase
	Names() []string
	LabelKeys(string) []string
}

func NewImageClient(client ImageClientBase) ImageClient {
	return &imageClient{
		ImageClientBase: client,
	}
}

type imageClient struct {
	ImageClientBase
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
			name = strconv.Itoa(img.ID)
		}
		names[i] = name
	}
	return names
}

// ImageLabelKeys returns a slice containing the keys of all labels assigned to
// the Image with the passed idOrName.
func (c *imageClient) LabelKeys(idOrName string) []string {
	img, _, err := c.Get(context.Background(), idOrName)
	if err != nil || img == nil || len(img.Labels) == 0 {
		return nil
	}
	return labelKeys(img.Labels)
}
