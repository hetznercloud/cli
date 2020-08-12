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
