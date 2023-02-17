package hcapi2

import (
	"context"
	"strconv"
)

// LocationClient embeds the Hetzner Cloud Location client and provides some
// additional helper functions.
type LocationClient interface {
	LocationClientBase
	Names() []string
	NetworkZones() []string
}

func NewLocationClient(client LocationClientBase) LocationClient {
	return &locationClient{
		LocationClientBase: client,
	}
}

type locationClient struct {
	LocationClientBase
}

// Names obtains a list of available locations. It returns nil if
// location names could not be fetched.
func (c *locationClient) Names() []string {
	locs, err := c.All(context.Background())
	if err != nil || len(locs) == 0 {
		return nil
	}
	names := make([]string, len(locs))
	for i, loc := range locs {
		name := loc.Name
		if name == "" {
			name = strconv.Itoa(loc.ID)
		}
		names[i] = name
	}
	return names
}

// NetworkZones obtains a list of available network zones. It returns nil if
// location data could not be fetched.
func (c *locationClient) NetworkZones() []string {
	locs, err := c.All(context.Background())
	if err != nil || len(locs) == 0 {
		return nil
	}

	zones := make(map[string]bool)
	for _, loc := range locs {
		if loc.NetworkZone != "" {
			zones[string(loc.NetworkZone)] = true
		}
	}

	var zoneList []string
	for zone := range zones {
		zoneList = append(zoneList, zone)
	}
	return zoneList
}
