package hcapi2

import (
	"context"
	"strconv"
)

// VolumeClient embeds the Hetzner Cloud Volume client and provides some additional
// helper functions.
type VolumeClient interface {
	VolumeClientBase
	Names() []string
}

func NewVolumeClient(client VolumeClientBase) VolumeClient {
	return &volumeClient{
		VolumeClientBase: client,
	}
}

type volumeClient struct {
	VolumeClientBase
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
			name = strconv.Itoa(vol.ID)
		}
		names[i] = name
	}
	return names
}

// VolumeLabelKeys returns a slice containing the keys of all labels assigned
// to the Volume with the passed idOrName.
func (c *volumeClient) VolumeLabelKeys(idOrName string) []string {
	vol, _, err := c.Get(context.Background(), idOrName)
	if err != nil || vol == nil || len(vol.Labels) == 0 {
		return nil
	}
	return labelKeys(vol.Labels)
}
