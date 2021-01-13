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
