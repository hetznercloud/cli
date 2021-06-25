package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type LoadBalancerTypeClient interface {
	LoadBalancerTypeClientBase
	Names() []string
	LoadBalancerTypeName(id int) string
	LoadBalancerTypeDescription(id int) string
}

func NewLoadBalancerTypeClient(client LoadBalancerTypeClientBase) LoadBalancerTypeClient {
	return &loadBalancerTypeClient{
		LoadBalancerTypeClientBase: client,
	}
}

type loadBalancerTypeClient struct {
	LoadBalancerTypeClientBase

	lbTypeByID map[int]*hcloud.LoadBalancerType
	once       sync.Once
	err        error
}

// LoadBalancerTypeName obtains the name of the loadBalancer type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *loadBalancerTypeClient) LoadBalancerTypeName(id int) string {
	if err := c.init(); err != nil {
		return strconv.Itoa(id)
	}

	loadBalancerType, ok := c.lbTypeByID[id]
	if !ok || loadBalancerType.Name == "" {
		return strconv.Itoa(id)
	}
	return loadBalancerType.Name
}

// LoadBalancerTypeDescription obtains the description of the loadBalancer type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *loadBalancerTypeClient) LoadBalancerTypeDescription(id int) string {
	if err := c.init(); err != nil {
		return strconv.Itoa(id)
	}

	loadBalancerType, ok := c.lbTypeByID[id]
	if !ok || loadBalancerType.Description == "" {
		return strconv.Itoa(id)
	}
	return loadBalancerType.Description
}

// Names returns a slice of all available loadBalancer types.
func (c *loadBalancerTypeClient) Names() []string {
	sts, err := c.All(context.Background())
	if err != nil || len(sts) == 0 {
		return nil
	}
	names := make([]string, len(sts))
	for i, st := range sts {
		names[i] = st.Name
	}
	return names
}

func (c *loadBalancerTypeClient) init() error {
	c.once.Do(func() {
		loadBalancerTypes, err := c.All(context.Background())
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(loadBalancerTypes) == 0 {
			return
		}
		c.lbTypeByID = make(map[int]*hcloud.LoadBalancerType, len(loadBalancerTypes))
		for _, srv := range loadBalancerTypes {
			c.lbTypeByID[srv.ID] = srv
		}
	})
	return c.err
}
