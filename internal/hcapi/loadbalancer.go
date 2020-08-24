package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// LoadBalancerClient embeds the Hetzner Cloud Load Balancer client and
// provides some additional helper functions. As a convenience it references
// the type client as well and provides a unified interface to Load Balancers
// and their types.
type LoadBalancerClient struct {
	*hcloud.LoadBalancerClient

	TypeClient *hcloud.LoadBalancerTypeClient
}

// LoadBalancerNames obtains a list of all available Load Balancer names. It
// returns nil if the names could not be fetched.
func (c *LoadBalancerClient) LoadBalancerNames() []string {
	lbs, err := c.All(context.Background())
	if err != nil || len(lbs) == 0 {
		return nil
	}
	names := make([]string, len(lbs))
	for i, lb := range lbs {
		name := lb.Name
		if name == "" {
			name = strconv.Itoa(lb.ID)
		}
		names[i] = name
	}
	return names
}

// LoadBalancerLabelKeys returns a slice containing the keys of all labels
// assigned to the Load Balancer with the passed idOrName.
func (c *LoadBalancerClient) LoadBalancerLabelKeys(idOrName string) []string {
	lb, _, err := c.Get(context.Background(), idOrName)
	if err != nil || lb == nil || len(lb.Labels) == 0 {
		return nil
	}
	return lkeys(lb.Labels)
}

// LoadBalancerTypeNames returns a slice containing the names of all available
// Load Balancer types.
func (c *LoadBalancerClient) LoadBalancerTypeNames() []string {
	types, err := c.TypeClient.All(context.Background())
	if err != nil || len(types) == 0 {
		return nil
	}
	names := make([]string, len(types))
	for i, typ := range types {
		names[i] = typ.Name
	}
	return names
}
