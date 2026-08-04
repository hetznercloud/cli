package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// LocationClient embeds the Hetzner Cloud Location client and provides some
// additional helper functions.
type LocationClient interface {
	hcloud.ILocationClient
	Names(context.Context) ([]string, error)
	NetworkZones(context.Context) ([]string, error)
}

func NewLocationClient(client hcloud.ILocationClient) LocationClient {
	return &locationClient{
		ILocationClient: client,
	}
}

type locationClient struct {
	hcloud.ILocationClient
}

// Names obtains a list of available locations. It returns nil if
// location names could not be fetched.
func (c *locationClient) Names(ctx context.Context) ([]string, error) {
	locations, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(locations, func(location *hcloud.Location) int64 { return location.ID }, func(location *hcloud.Location) string { return location.Name }), nil
}

// NetworkZones obtains a list of available network zones. It returns nil if
// location data could not be fetched.
func (c *locationClient) NetworkZones(ctx context.Context) ([]string, error) {
	locations, err := c.All(ctx)
	if err != nil {
		return nil, err
	}

	zones := make(map[string]bool)
	for _, loc := range locations {
		if loc.NetworkZone != "" {
			zones[string(loc.NetworkZone)] = true
		}
	}

	var zoneList []string
	for zone := range zones {
		zoneList = append(zoneList, zone)
	}
	return zoneList, nil
}
