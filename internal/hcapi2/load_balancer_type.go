package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type LoadBalancerTypeClient interface {
	hcloud.ILoadBalancerTypeClient
	Names(context.Context) ([]string, error)
}

func NewLoadBalancerTypeClient(client hcloud.ILoadBalancerTypeClient) LoadBalancerTypeClient {
	return &loadBalancerTypeClient{
		ILoadBalancerTypeClient: client,
	}
}

type loadBalancerTypeClient struct {
	hcloud.ILoadBalancerTypeClient
}

// Names returns a slice of all available loadBalancer types.
func (c *loadBalancerTypeClient) Names(ctx context.Context) ([]string, error) {
	loadBalancerTypes, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(loadBalancerTypes))
	for i, st := range loadBalancerTypes {
		names[i] = st.Name
	}
	return names, nil
}
