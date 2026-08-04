package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// VolumeClient embeds the Hetzner Cloud Volume client and provides some additional
// helper functions.
type VolumeClient interface {
	hcloud.IVolumeClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
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
func (c *volumeClient) Names(ctx context.Context) ([]string, error) {
	volumes, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(volumes, func(volume *hcloud.Volume) int64 { return volume.ID }, func(volume *hcloud.Volume) string { return volume.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels assigned
// to the Volume with the passed idOrName.
func (c *volumeClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	volume, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if volume == nil {
		return nil, nil
	}
	return labelKeys(volume.Labels), nil
}
