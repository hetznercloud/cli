package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// LoadBalancerClient embeds the Hetzner Cloud LoadBalancer client and provides some
// additional helper functions.
type LoadBalancerClient interface {
	hcloud.ILoadBalancerClient
	LoadBalancerName(id int64) string
	Names() []string
	LabelKeys(string) []string
}

func NewLoadBalancerClient(client hcloud.ILoadBalancerClient) LoadBalancerClient {
	return &loadBalancerClient{
		ILoadBalancerClient: client,
	}
}

type loadBalancerClient struct {
	hcloud.ILoadBalancerClient

	lbByID map[int64]*hcloud.LoadBalancer

	once sync.Once
	err  error
}

// LoadBalancerName obtains the name of the server with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *loadBalancerClient) LoadBalancerName(id int64) string {
	if err := c.init(); err != nil {
		return strconv.FormatInt(id, 10)
	}

	lb, ok := c.lbByID[id]
	if !ok || lb.Name == "" {
		return strconv.FormatInt(id, 10)
	}
	return lb.Name
}

// Names obtains a list of available data centers. It returns nil if
// data center names could not be fetched.
func (c *loadBalancerClient) Names() []string {
	dcs, err := c.All(context.Background())
	if err != nil || len(dcs) == 0 {
		return nil
	}
	names := make([]string, len(dcs))
	for i, dc := range dcs {
		name := dc.Name
		if name == "" {
			name = strconv.FormatInt(dc.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the loadBalancer with the passed idOrName.
func (c *loadBalancerClient) LabelKeys(idOrName string) []string {
	loadBalancer, _, err := c.Get(context.Background(), idOrName)
	if err != nil || loadBalancer == nil || len(loadBalancer.Labels) == 0 {
		return nil
	}
	return labelKeys(loadBalancer.Labels)
}

func (c *loadBalancerClient) init() error {
	c.once.Do(func() {
		srvs, err := c.All(context.Background())
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(srvs) == 0 {
			return
		}
		c.lbByID = make(map[int64]*hcloud.LoadBalancer, len(srvs))
		for _, srv := range srvs {
			c.lbByID[srv.ID] = srv
		}
	})
	return c.err
}
