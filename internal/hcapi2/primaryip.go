package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// PrimaryIPClient embeds the hcloud PrimaryIPClient (via an interface) and provides
// some additional helper functions.
type PrimaryIPClient interface {
	hcloud.IPrimaryIPClient
	Names(hideAssigned, hideUnassigned bool, ipType *hcloud.PrimaryIPType) func() []string
	LabelKeys(idOrName string) []string
}

// NewPrimaryIPClient creates a new primary IP client.
func NewPrimaryIPClient(client hcloud.IPrimaryIPClient) PrimaryIPClient {
	return &primaryIPClient{
		IPrimaryIPClient: client,
	}
}

// PrimaryIPClient embeds the Hetzner Cloud PrimaryIP client and provides some
// additional helper functions.
type primaryIPClient struct {
	hcloud.IPrimaryIPClient
}

// Names obtains a list of available primary IPs. It returns nil if
// no primary IP names could be fetched or none were available.
// hideUnassigned: if true, only returns names of primary IPs that are assigned to a server
// hideAssigned: if true, only returns names of primary IPs that are not assigned to a server
// ipType: if not nil, only returns primary IPs of the specified type (IPv4 or IPv6)
// Returns a func() []string so that the list can be lazily evaluated
func (c *primaryIPClient) Names(hideAssigned, hideUnassigned bool, ipType *hcloud.PrimaryIPType) func() []string {
	return func() []string {
		fips, err := c.All(context.Background())
		if err != nil || len(fips) == 0 {
			return nil
		}
		names := make([]string, len(fips))
		for i, fip := range fips {
			if (hideAssigned && fip.AssigneeID > 0) ||
				(hideUnassigned && fip.AssigneeID == 0) ||
				(ipType != nil && fip.Type != *ipType) {
				continue
			}
			name := fip.Name
			if name == "" {
				name = strconv.FormatInt(fip.ID, 10)
			}
			names[i] = name
		}
		return names
	}
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the Primary IP with the passed idOrName.
func (c *primaryIPClient) LabelKeys(idOrName string) []string {
	fip, _, err := c.Get(context.Background(), idOrName)
	if err != nil || fip == nil || len(fip.Labels) == 0 {
		return nil
	}
	return labelKeys(fip.Labels)
}
