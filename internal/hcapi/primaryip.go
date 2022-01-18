package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// PrimaryIPClient embeds the Hetzner Cloud PrimaryIP client and provides some
// additional helper functions.
type PrimaryIPClient struct {
	*hcloud.PrimaryIPClient
}

// PrimaryIPNames obtains a list of available floating IPs. It returns nil if
// no floating IP names could be fetched or none were available.
func (c *PrimaryIPClient) PrimaryIPNames() []string {
	primaryIPs, err := c.All(context.Background())
	if err != nil || len(primaryIPs) == 0 {
		return nil
	}
	names := make([]string, len(primaryIPs))
	for i, primaryIP := range primaryIPs {
		name := primaryIP.Name
		if name == "" {
			name = strconv.Itoa(primaryIP.ID)
		}
		names[i] = name
	}
	return names
}

// PrimaryIPLabelKeys returns a slice containing the keys of all labels
// assigned to the Primary IP with the passed idOrName.
func (c *PrimaryIPClient) PrimaryIPLabelKeys(idOrName string) []string {
	primaryIP, _, err := c.Get(context.Background(), idOrName)
	if err != nil || primaryIP == nil || len(primaryIP.Labels) == 0 {
		return nil
	}
	return lkeys(primaryIP.Labels)
}
