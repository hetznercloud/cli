package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// VolumeClient embeds the Hetzner Cloud Volume client and provides some additional
// helper functions.
type VolumeClient interface {
	hcloud.IVolumeClient
	Names() []string
	LabelKeys(idOrName string) []string
}

func NewVolumeClient(client hcloud.IVolumeClient) VolumeClient {
	return &volumeClient{
		IVolumeClient: client,
	}
}

type volumeClient struct {
	hcloud.IVolumeClient
}

// Names obtains a list of available volumes for the current account. It
// returns nil if the current project has no volumes or the volume names could
// not be fetched.
func (c *volumeClient) Names() []string {
	vols, err := c.All(context.Background())
	if err != nil || len(vols) == 0 {
		return nil
	}
	names := make([]string, len(vols))
	for i, vol := range vols {
		name := vol.Name
		if name == "" {
			name = strconv.FormatInt(vol.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels assigned
// to the Volume with the passed idOrName.
func (c *volumeClient) LabelKeys(idOrName string) []string {
	vol, _, err := c.Get(context.Background(), idOrName)
	if err != nil || vol == nil || len(vol.Labels) == 0 {
		return nil
	}
	return labelKeys(vol.Labels)
}
