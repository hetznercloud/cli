package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// FloatingIPClient embeds the hcloud FloatingIPClient (via an interface) and provides
// some additional helper functions.
type FloatingIPClient interface {
	hcloud.IFloatingIPClient
	Names() []string
	LabelKeys(idOrName string) []string
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
func (c *floatingIPClient) Names() []string {
	fips, err := c.All(context.Background())
	if err != nil || len(fips) == 0 {
		return nil
	}
	names := make([]string, len(fips))
	for i, fip := range fips {
		name := fip.Name
		if name == "" {
			name = strconv.FormatInt(fip.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the Floating IP with the passed idOrName.
func (c *floatingIPClient) LabelKeys(idOrName string) []string {
	fip, _, err := c.Get(context.Background(), idOrName)
	if err != nil || fip == nil || len(fip.Labels) == 0 {
		return nil
	}
	return labelKeys(fip.Labels)
}
