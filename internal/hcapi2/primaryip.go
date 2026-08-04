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
	Names(hideAssigned, hideUnassigned bool, ipType *hcloud.PrimaryIPType) func(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
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
// Returns a function so that the list can be lazily evaluated with the command context.
func (c *primaryIPClient) Names(hideAssigned, hideUnassigned bool, ipType *hcloud.PrimaryIPType) func(context.Context) ([]string, error) {
	return func(ctx context.Context) ([]string, error) {
		primaryIPs, err := c.All(ctx)
		if err != nil {
			return nil, err
		}
		names := make([]string, 0, len(primaryIPs))
		for _, primaryIP := range primaryIPs {
			if (hideAssigned && primaryIP.AssigneeID > 0) ||
				(hideUnassigned && primaryIP.AssigneeID == 0) ||
				(ipType != nil && primaryIP.Type != *ipType) {
				continue
			}
			name := primaryIP.Name
			if name == "" {
				name = strconv.FormatInt(primaryIP.ID, 10)
			}
			names = append(names, name)
		}
		return names, nil
	}
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the Primary IP with the passed idOrName.
func (c *primaryIPClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	primaryIP, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if primaryIP == nil {
		return nil, nil
	}
	return labelKeys(primaryIP.Labels), nil
}
