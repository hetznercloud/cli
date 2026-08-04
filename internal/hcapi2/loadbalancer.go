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
	LoadBalancerName(context.Context, int64) (string, error)
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
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

// LoadBalancerName obtains the name of the load balancer with id. It returns
// the numeric ID when the API response does not contain a matching named load balancer.
func (c *loadBalancerClient) LoadBalancerName(ctx context.Context, id int64) (string, error) {
	if err := c.init(ctx); err != nil {
		return "", err
	}

	lb, ok := c.lbByID[id]
	if !ok || lb.Name == "" {
		return strconv.FormatInt(id, 10), nil
	}
	return lb.Name, nil
}

// Names obtains a list of available data centers. It returns nil if
// data center names could not be fetched.
func (c *loadBalancerClient) Names(ctx context.Context) ([]string, error) {
	loadBalancers, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(loadBalancers, func(loadBalancer *hcloud.LoadBalancer) int64 { return loadBalancer.ID }, func(loadBalancer *hcloud.LoadBalancer) string { return loadBalancer.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the loadBalancer with the passed idOrName.
func (c *loadBalancerClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	loadBalancer, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if loadBalancer == nil {
		return nil, nil
	}
	return labelKeys(loadBalancer.Labels), nil
}

func (c *loadBalancerClient) init(ctx context.Context) error {
	c.once.Do(func() {
		srvs, err := c.All(ctx)
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
