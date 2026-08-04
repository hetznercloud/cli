package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DatacenterClient embeds the Hetzner Cloud DataCenter client and provides some
// additional helper functions.
type DatacenterClient interface {
	hcloud.IDatacenterClient
	Names(context.Context) ([]string, error)
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
func (c *datacenterClient) Names(ctx context.Context) ([]string, error) {
	datacenters, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(datacenters, func(datacenter *hcloud.Datacenter) int64 { return datacenter.ID }, func(datacenter *hcloud.Datacenter) string { return datacenter.Name }), nil
}
