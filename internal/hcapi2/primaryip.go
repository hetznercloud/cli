package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// PrimaryIPClient embeds the hcloud PrimaryIPClient (via an interface) and provides
// some additional helper functions.
type PrimaryIPClient interface {
	PrimaryIPClientBase
	Names() []string
	IPv4Names() []string
	IPv6Names() []string
	LabelKeys(idOrName string) []string
}

// NewPrimaryIPClient creates a new primary IP client.
func NewPrimaryIPClient(client PrimaryIPClientBase) PrimaryIPClient {
	return &primaryIPClient{
		PrimaryIPClientBase: client,
	}
}

// PrimaryIPClient embeds the Hetzner Cloud PrimaryIP client and provides some
// additional helper functions.
type primaryIPClient struct {
	PrimaryIPClientBase
}

// Names obtains a list of available primary IPs. It returns nil if
// no primary IP names could be fetched or none were available.
func (c *primaryIPClient) Names() []string {
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

// IPv4Names obtains a list of available primary IPv4s. It returns nil if
// no primary IP names could be fetched or none were available.
func (c *primaryIPClient) IPv4Names() []string {
	fips, err := c.All(context.Background())
	if err != nil || len(fips) == 0 {
		return nil
	}
	names := []string{}
	for _, fip := range fips {
		if fip.Type == hcloud.PrimaryIPTypeIPv4 {
			name := fip.Name
			if name == "" {
				name = strconv.Itoa(fip.ID)
			}
			names = append(names, name)
		}
	}
	return names
}

// IPv6Names obtains a list of available primary IPv6s. It returns nil if
// no primary IP names could be fetched or none were available.
func (c *primaryIPClient) IPv6Names() []string {
	fips, err := c.All(context.Background())
	if err != nil || len(fips) == 0 {
		return nil
	}
	names := []string{}
	for _, fip := range fips {
		if fip.Type == hcloud.PrimaryIPTypeIPv6 {
			name := fip.Name
			if name == "" {
				name = strconv.Itoa(fip.ID)
			}
			names = append(names, name)
		}
	}
	return names
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
