package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// NetworkClient embeds the Hetzner Cloud Network client and provides some
// additional helper functions.
type NetworkClient interface {
	NetworkClientBase
	Names() []string
	Name(int) string
	LabelKeys(string) []string
}

func NewNetworkClient(client NetworkClientBase) NetworkClient {
	return &networkClient{
		NetworkClientBase: client,
	}
}

type networkClient struct {
	NetworkClientBase

	netsByID   map[int]*hcloud.Network
	netsByName map[string]*hcloud.Network

	once sync.Once
	err  error
}

// NetworkName obtains the name of the network with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *networkClient) Name(id int) string {
	if err := c.init(); err != nil {
		return strconv.Itoa(id)
	}

	net, ok := c.netsByID[id]
	if !ok || net.Name == "" {
		return strconv.Itoa(id)
	}
	return net.Name
}

// Names obtains a list of available networks. It returns nil if the
// network names could not be fetched or if there are no networks.
func (c *networkClient) Names() []string {
	if err := c.init(); err != nil || len(c.netsByID) == 0 {
		return nil
	}
	names := make([]string, len(c.netsByID))
	i := 0
	for _, net := range c.netsByID {
		name := net.Name
		if name == "" {
			name = strconv.Itoa(net.ID)
		}
		names[i] = name
		i++
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels assigned
// to the Network with the passed idOrName.
func (c *networkClient) LabelKeys(idOrName string) []string {
	var net *hcloud.Network

	if err := c.init(); err != nil || len(c.netsByID) == 0 {
		return nil
	}
	if id, err := strconv.Atoi(idOrName); err != nil {
		net = c.netsByID[id]
	}
	if v, ok := c.netsByName[idOrName]; ok && net == nil {
		net = v
	}
	if net == nil || len(net.Labels) == 0 {
		return nil
	}
	return labelKeys(net.Labels)
}

func (c *networkClient) init() error {
	c.once.Do(func() {
		nets, err := c.All(context.Background())
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(nets) == 0 {
			return
		}
		c.netsByID = make(map[int]*hcloud.Network, len(nets))
		c.netsByName = make(map[string]*hcloud.Network, len(nets))
		for _, net := range nets {
			c.netsByID[net.ID] = net
			c.netsByName[net.Name] = net
		}
	})
	return c.err
}
