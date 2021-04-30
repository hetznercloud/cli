package hcapi2

import (
	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// Client makes all API clients accessible via a single interface.
type Client interface {
	Datacenter() DatacenterClient
	Firewall() FirewallClient
	FloatingIP() FloatingIPClient
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

// NewClient creates a new CLI API client extending hcloud.Client.
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

func (c *client) FloatingIP() FloatingIPClient {
	return NewFloatingIPClient(&c.client.FloatingIP)
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

type MockClient struct {
	DatacenterClient *MockDatacenterClient
	FirewallClient   *MockFirewallClient
	FloatingIPClient *MockFloatingIPClient
	ImageClient      *MockImageClient
	LocationClient   *MockLocationClient
	NetworkClient    *MockNetworkClient
	ServerClient     *MockServerClient
	ServerTypeClient *MockServerTypeClient
	SSHKeyClient     *MockSSHKeyClient
	VolumeClient     *MockVolumeClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	return &MockClient{
		DatacenterClient: NewMockDatacenterClient(ctrl),
		FirewallClient:   NewMockFirewallClient(ctrl),
		FloatingIPClient: NewMockFloatingIPClient(ctrl),
		ImageClient:      NewMockImageClient(ctrl),
		LocationClient:   NewMockLocationClient(ctrl),
		NetworkClient:    NewMockNetworkClient(ctrl),
		ServerClient:     NewMockServerClient(ctrl),
		ServerTypeClient: NewMockServerTypeClient(ctrl),
		SSHKeyClient:     NewMockSSHKeyClient(ctrl),
		VolumeClient:     NewMockVolumeClient(ctrl),
	}
}

func (c *MockClient) Datacenter() DatacenterClient {
	return c.DatacenterClient
}

func (c *MockClient) Firewall() FirewallClient {
	return c.FirewallClient
}

func (c *MockClient) FloatingIP() FloatingIPClient {
	return c.FloatingIPClient
}

func (c *MockClient) Image() ImageClient {
	return c.ImageClient
}

func (c *MockClient) Location() LocationClient {
	return c.LocationClient
}

func (c *MockClient) Network() NetworkClient {
	return c.NetworkClient
}

func (c *MockClient) Server() ServerClient {
	return c.ServerClient
}

func (c *MockClient) ServerType() ServerTypeClient {
	return c.ServerTypeClient
}

func (c *MockClient) SSHKey() SSHKeyClient {
	return c.SSHKeyClient
}

func (c *MockClient) Volume() VolumeClient {
	return c.VolumeClient
}
