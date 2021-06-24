package hcapi2

import (
	"sync"

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
	Certificate() CertificateClient
	LoadBalancer() LoadBalancerClient
	ISO() ISOClient
}

type client struct {
	client             *hcloud.Client
	certificateClient  CertificateClient
	datacenterClient   DatacenterClient
	serverClient       ServerClient
	serverTypeClient   ServerTypeClient
	locationClient     LocationClient
	loadBalancerClient LoadBalancerClient
	networkClient      NetworkClient
	firewallClient     FirewallClient
	floatingIPClient   FloatingIPClient
	imageClient        ImageClient
	isoClient          ISOClient
	sshKeyClient       SSHKeyClient
	volumeClient       VolumeClient

	mu sync.Mutex
}

// NewClient creates a new CLI API client extending hcloud.Client.
func NewClient(c *hcloud.Client) Client {
	return &client{
		client: c,
	}
}
func (c *client) Certificate() CertificateClient {
	c.mu.Lock()
	if c.certificateClient == nil {
		c.certificateClient = NewCertificateClient(&c.client.Certificate)
	}
	defer c.mu.Unlock()
	return c.certificateClient
}

func (c *client) Datacenter() DatacenterClient {
	c.mu.Lock()
	if c.datacenterClient == nil {
		c.datacenterClient = NewDatacenterClient(&c.client.Datacenter)
	}
	defer c.mu.Unlock()
	return c.datacenterClient
}

func (c *client) Firewall() FirewallClient {
	c.mu.Lock()
	if c.firewallClient == nil {
		c.firewallClient = NewFirewallClient(&c.client.Firewall)
	}
	defer c.mu.Unlock()
	return c.firewallClient
}

func (c *client) FloatingIP() FloatingIPClient {
	c.mu.Lock()
	if c.floatingIPClient == nil {
		c.floatingIPClient = NewFloatingIPClient(&c.client.FloatingIP)
	}
	defer c.mu.Unlock()
	return c.floatingIPClient
}

func (c *client) Image() ImageClient {
	c.mu.Lock()
	if c.imageClient == nil {
		c.imageClient = NewImageClient(&c.client.Image)
	}
	defer c.mu.Unlock()
	return c.imageClient
}

func (c *client) ISO() ISOClient {
	c.mu.Lock()
	if c.isoClient == nil {
		c.isoClient = NewISOClient(&c.client.ISO)
	}
	defer c.mu.Unlock()
	return c.isoClient
}

func (c *client) Location() LocationClient {
	c.mu.Lock()
	if c.locationClient == nil {
		c.locationClient = NewLocationClient(&c.client.Location)
	}
	defer c.mu.Unlock()
	return c.locationClient
}

func (c *client) LoadBalancer() LoadBalancerClient {
	c.mu.Lock()
	if c.loadBalancerClient == nil {
		c.loadBalancerClient = NewLoadBalancerClient(&c.client.LoadBalancer)
	}
	defer c.mu.Unlock()
	return c.loadBalancerClient
}

func (c *client) Network() NetworkClient {
	c.mu.Lock()
	if c.networkClient == nil {
		c.networkClient = NewNetworkClient(&c.client.Network)
	}
	defer c.mu.Unlock()
	return c.networkClient
}

func (c *client) Server() ServerClient {
	c.mu.Lock()
	if c.serverClient == nil {
		c.serverClient = NewServerClient(&c.client.Server)
	}
	defer c.mu.Unlock()
	return c.serverClient
}

func (c *client) ServerType() ServerTypeClient {
	c.mu.Lock()
	if c.serverTypeClient == nil {
		c.serverTypeClient = NewServerTypeClient(&c.client.ServerType)
	}
	defer c.mu.Unlock()
	return c.serverTypeClient
}

func (c *client) SSHKey() SSHKeyClient {
	c.mu.Lock()
	if c.sshKeyClient == nil {
		c.sshKeyClient = NewSSHKeyClient(&c.client.SSHKey)
	}
	defer c.mu.Unlock()
	return c.sshKeyClient
}

func (c *client) Volume() VolumeClient {
	c.mu.Lock()
	if c.volumeClient == nil {
		c.volumeClient = NewVolumeClient(&c.client.Volume)
	}
	defer c.mu.Unlock()
	return c.volumeClient
}

type MockClient struct {
	CertificateClient  *MockCertificateClient
	DatacenterClient   *MockDatacenterClient
	FirewallClient     *MockFirewallClient
	FloatingIPClient   *MockFloatingIPClient
	ImageClient        *MockImageClient
	LocationClient     *MockLocationClient
	LoadBalancerClient *MockLoadBalancerClient
	NetworkClient      *MockNetworkClient
	ServerClient       *MockServerClient
	ServerTypeClient   *MockServerTypeClient
	SSHKeyClient       *MockSSHKeyClient
	VolumeClient       *MockVolumeClient
	ISOClient          *MockISOClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	return &MockClient{
		CertificateClient:  NewMockCertificateClient(ctrl),
		DatacenterClient:   NewMockDatacenterClient(ctrl),
		FirewallClient:     NewMockFirewallClient(ctrl),
		FloatingIPClient:   NewMockFloatingIPClient(ctrl),
		ImageClient:        NewMockImageClient(ctrl),
		LocationClient:     NewMockLocationClient(ctrl),
		LoadBalancerClient: NewMockLoadBalancerClient(ctrl),
		NetworkClient:      NewMockNetworkClient(ctrl),
		ServerClient:       NewMockServerClient(ctrl),
		ServerTypeClient:   NewMockServerTypeClient(ctrl),
		SSHKeyClient:       NewMockSSHKeyClient(ctrl),
		VolumeClient:       NewMockVolumeClient(ctrl),
		ISOClient:          NewMockISOClient(ctrl),
	}
}
func (c *MockClient) Certificate() CertificateClient {
	return c.CertificateClient
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

func (c *MockClient) ISO() ISOClient {
	return c.ISOClient
}

func (c *MockClient) Location() LocationClient {
	return c.LocationClient
}

func (c *MockClient) LoadBalancer() LoadBalancerClient {
	return c.LoadBalancerClient
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
