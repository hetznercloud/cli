package hcapi

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// LoadBalancerClient embeds the Hetzner Cloud Load Balancer client and
// provides some additional helper functions. As a convenience it references
// the type client as well and provides a unified interface to Load Balancers
// and their types.
type LoadBalancerClient struct {
	*hcloud.LoadBalancerClient

	TypeClient *hcloud.LoadBalancerTypeClient

	lbByID   map[int]*hcloud.LoadBalancer
	lbByName map[string]*hcloud.LoadBalancer
	once     sync.Once
	err      error
}

// LoadBalancerNames obtains a list of all available Load Balancer names. It
// returns nil if the names could not be fetched.
func (c *LoadBalancerClient) LoadBalancerNames() []string {
	if err := c.init(); err != nil || len(c.lbByName) == 0 {
		return nil
	}
	names := make([]string, len(c.lbByName))
	i := 0
	for name, lb := range c.lbByName {
		if name == "" {
			name = strconv.Itoa(lb.ID)
		}
		names[i] = name
		i++
	}
	return names
}

// LoadBalancerName returns a Load Balancer's name given its id.
//
// If a Load Balancer has no name the Load Balancer's id is returned converted
// to a string.
func (c *LoadBalancerClient) LoadBalancerName(id int) string {
	if err := c.init(); err != nil || len(c.lbByName) == 0 {
		return ""
	}
	lb, ok := c.lbByID[id]
	if !ok {
		return ""
	}
	if lb.Name == "" {
		return strconv.Itoa(lb.ID)
	}
	return lb.Name
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

func (c *LoadBalancerClient) init() error {
	c.once.Do(func() {
		lbs, err := c.All(context.Background())
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(lbs) == 0 {
			return
		}
		c.lbByID = make(map[int]*hcloud.LoadBalancer, len(lbs))
		c.lbByName = make(map[string]*hcloud.LoadBalancer, len(lbs))
		for _, lb := range lbs {
			c.lbByID[lb.ID] = lb
			c.lbByName[lb.Name] = lb
		}
	})
	return c.err
}
