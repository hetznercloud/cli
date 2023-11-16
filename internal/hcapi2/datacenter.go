package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DatacenterClient embeds the Hetzner Cloud DataCenter client and provides some
// additional helper functions.
type DatacenterClient interface {
	hcloud.IDatacenterClient
	Names() []string
}

func NewDatacenterClient(client hcloud.IDatacenterClient) DatacenterClient {
	return &datacenterClient{
		IDatacenterClient: client,
	}
}

type datacenterClient struct {
	hcloud.IDatacenterClient
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
			name = strconv.FormatInt(dc.ID, 10)
		}
		names[i] = name
	}
	return names
}
