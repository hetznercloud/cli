package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// NetworkClient embeds the Hetzner Cloud Network client and provides some
// additional helper functions.
type NetworkClient interface {
	hcloud.INetworkClient
	Names(context.Context) ([]string, error)
	Name(context.Context, int64) (string, error)
	LabelKeys(context.Context, string) ([]string, error)
}

func NewNetworkClient(client hcloud.INetworkClient) NetworkClient {
	return &networkClient{
		INetworkClient: client,
	}
}

type networkClient struct {
	hcloud.INetworkClient

	netsByID   map[int64]*hcloud.Network
	netsByName map[string]*hcloud.Network

	once sync.Once
	err  error
}

// Name obtains the name of the network with id. It returns the numeric ID when
// the API response does not contain a matching named network.
func (c *networkClient) Name(ctx context.Context, id int64) (string, error) {
	if err := c.init(ctx); err != nil {
		return "", err
	}

	net, ok := c.netsByID[id]
	if !ok || net.Name == "" {
		return strconv.FormatInt(id, 10), nil
	}
	return net.Name, nil
}

// Names obtains a list of available networks. It returns nil if the
// network names could not be fetched or if there are no networks.
func (c *networkClient) Names(ctx context.Context) ([]string, error) {
	networks, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(networks, func(network *hcloud.Network) int64 { return network.ID }, func(network *hcloud.Network) string { return network.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels assigned
// to the Network with the passed idOrName.
func (c *networkClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	network, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if network == nil {
		return nil, nil
	}
	return labelKeys(network.Labels), nil
}

func (c *networkClient) init(ctx context.Context) error {
	c.once.Do(func() {
		nets, err := c.All(ctx)
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(nets) == 0 {
			return
		}
		c.netsByID = make(map[int64]*hcloud.Network, len(nets))
		c.netsByName = make(map[string]*hcloud.Network, len(nets))
		for _, net := range nets {
			c.netsByID[net.ID] = net
			c.netsByName[net.Name] = net
		}
	})
	return c.err
}
