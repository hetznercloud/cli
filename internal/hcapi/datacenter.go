package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// DataCenterClient embeds the Hetzner Cloud DataCenter client and provides some
// additional helper functions.
type DataCenterClient struct {
	*hcloud.DatacenterClient
}

// DataCenterNames obtains a list of available data centers. It returns nil if
// data center names could not be fetched.
func (c *DataCenterClient) DataCenterNames() []string {
	dcs, err := c.All(context.Background())
	if err != nil || len(dcs) == 0 {
		return nil
	}
	names := make([]string, len(dcs))
	for i, dc := range dcs {
		name := dc.Name
		if name == "" {
			name = strconv.Itoa(dc.ID)
		}
		names[i] = name
	}
	return names
}
