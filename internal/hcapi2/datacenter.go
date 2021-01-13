package hcapi2

import (
	"context"
	"strconv"
)

// DatacenterClient embeds the Hetzner Cloud DataCenter client and provides some
// additional helper functions.
type DatacenterClient interface {
	DatacenterClientBase
	Names() []string
}

func NewDatacenterClient(client DatacenterClientBase) DatacenterClient {
	return &datacenterClient{
		DatacenterClientBase: client,
	}
}

type datacenterClient struct {
	DatacenterClientBase
}

// Names obtains a list of available data centers. It returns nil if
// data center names could not be fetched.
func (c *datacenterClient) Names() []string {
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
