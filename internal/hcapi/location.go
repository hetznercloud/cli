package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// LocationClient embeds the Hetzner Cloud Location client and provides some
// additional helper functions.
type LocationClient struct {
	*hcloud.LocationClient
}

// LocationNames obtains a list of available locations. It returns nil if
// location names could not be fetched.
func (c *LocationClient) LocationNames() []string {
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

// NetworkZoneNames obtains a list of available network zones. It returns nil if
// network zone names could not be fetched.
func (c *LocationClient) NetworkZoneNames() []string {
	locs, err := c.All(context.Background())
	if err != nil || len(locs) == 0 {
		return nil
	}
	// Use map to get unique elements
	namesMap := map[hcloud.NetworkZone]bool{}
	for _, loc := range locs {
		name := loc.NetworkZone
		namesMap[name] = true
	}

	// Unique names from map to slice
	names := make([]string, len(namesMap))
	for name := range namesMap {
		names = append(names, string(name))
	}

	return names
}
