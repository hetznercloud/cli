package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// FloatingIPClient embeds the hcloud FloatingIPClient (via an interface) and provides
// some additional helper functions.
type FloatingIPClient interface {
	hcloud.IFloatingIPClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
}

// NewFloatingIPClient creates a new floating IP client.
func NewFloatingIPClient(client hcloud.IFloatingIPClient) FloatingIPClient {
	return &floatingIPClient{
		IFloatingIPClient: client,
	}
}

// FloatingIPClient embeds the Hetzner Cloud FloatingIP client and provides some
// additional helper functions.
type floatingIPClient struct {
	hcloud.IFloatingIPClient
}

// Names obtains a list of available floating IPs. It returns nil if
// no floating IP names could be fetched or none were available.
func (c *floatingIPClient) Names(ctx context.Context) ([]string, error) {
	fips, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(fips, func(fip *hcloud.FloatingIP) int64 { return fip.ID }, func(fip *hcloud.FloatingIP) string { return fip.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the Floating IP with the passed idOrName.
func (c *floatingIPClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	fip, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if fip == nil {
		return nil, nil
	}
	return labelKeys(fip.Labels), nil
}
