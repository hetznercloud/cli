package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// FloatingIPClient embeds the Hetzner Cloud FloatingIP client and provides some
// additional helper functions.
type FloatingIPClient struct {
	*hcloud.FloatingIPClient
}

// FloatingIPNames obtains a list of available floating IPs. It returns nil if
// no floating IP names could be fetched or none were available.
func (c *FloatingIPClient) FloatingIPNames() []string {
	fips, err := c.All(context.Background())
	if err != nil || len(fips) == 0 {
		return nil
	}
	names := make([]string, len(fips))
	for i, fip := range fips {
		name := fip.Name
		if name == "" {
			name = strconv.Itoa(fip.ID)
		}
		names[i] = name
	}
	return names
}

// FloatingIPLabelKeys returns a slice containing the keys of all labels
// assigned to the Floating IP with the passed idOrName.
func (c *FloatingIPClient) FloatingIPLabelKeys(idOrName string) []string {
	fip, _, err := c.Get(context.Background(), idOrName)
	if err != nil || fip == nil || len(fip.Labels) == 0 {
		return nil
	}
	return lkeys(fip.Labels)
}
