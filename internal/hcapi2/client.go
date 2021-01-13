package hcapi2

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Client interface {
	Datacenter() DatacenterClient
	Firewall() FirewallClient
	Image() ImageClient
	Location() LocationClient
	Network() NetworkClient
	Server() ServerClient
	ServerType() ServerTypeClient
	SSHKey() SSHKeyClient
	Volume() VolumeClient
}

type client struct {
	client *hcloud.Client
}

func NewClient(c *hcloud.Client) Client {
	return &client{
		client: c,
	}
}

func (c *client) Datacenter() DatacenterClient {
	return NewDatacenterClient(&c.client.Datacenter)
}

func (c *client) Firewall() FirewallClient {
	return NewFirewallClient(&c.client.Firewall)
}

func (c *client) Image() ImageClient {
	return NewImageClient(&c.client.Image)
}

func (c *client) Location() LocationClient {
	return NewLocationClient(&c.client.Location)
}

func (c *client) Network() NetworkClient {
	return NewNetworkClient(&c.client.Network)
}

func (c *client) Server() ServerClient {
	return NewServerClient(&c.client.Server)
}

func (c *client) ServerType() ServerTypeClient {
	return NewServerTypeClient(&c.client.ServerType)
}

func (c *client) SSHKey() SSHKeyClient {
	return NewSSHKeyClient(&c.client.SSHKey)
}

func (c *client) Volume() VolumeClient {
	return NewVolumeClient(&c.client.Volume)
}
