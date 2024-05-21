package hcapi2

import (
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mock"
)

type mockClient struct {
	client *mock.Client
	cache  clientCache
	mu     sync.Mutex
}

func NewMockClient(client *mock.Client) Client {
	c := &mockClient{
		client: client,
	}
	return c
}

func (*mockClient) WithOpts(...hcloud.ClientOption) {}

func (c *mockClient) Action() ActionClient {
	c.mu.Lock()
	if c.cache.actionClient == nil {
		c.cache.actionClient = NewActionClient(c.client.Action)
	}
	defer c.mu.Unlock()
	return c.cache.actionClient
}

func (c *mockClient) Certificate() CertificateClient {
	c.mu.Lock()
	if c.cache.certificateClient == nil {
		c.cache.certificateClient = NewCertificateClient(c.client.Certificate)
	}
	defer c.mu.Unlock()
	return c.cache.certificateClient
}

func (c *mockClient) Datacenter() DatacenterClient {
	c.mu.Lock()
	if c.cache.datacenterClient == nil {
		c.cache.datacenterClient = NewDatacenterClient(c.client.Datacenter)
	}
	defer c.mu.Unlock()
	return c.cache.datacenterClient
}

func (c *mockClient) Firewall() FirewallClient {
	c.mu.Lock()
	if c.cache.firewallClient == nil {
		c.cache.firewallClient = NewFirewallClient(c.client.Firewall)
	}
	defer c.mu.Unlock()
	return c.cache.firewallClient
}

func (c *mockClient) FloatingIP() FloatingIPClient {
	c.mu.Lock()
	if c.cache.floatingIPClient == nil {
		c.cache.floatingIPClient = NewFloatingIPClient(c.client.FloatingIP)
	}
	defer c.mu.Unlock()
	return c.cache.floatingIPClient
}

func (c *mockClient) PrimaryIP() PrimaryIPClient {
	c.mu.Lock()
	if c.cache.primaryIPClient == nil {
		c.cache.primaryIPClient = NewPrimaryIPClient(c.client.PrimaryIP)
	}
	defer c.mu.Unlock()
	return c.cache.primaryIPClient
}

func (c *mockClient) Image() ImageClient {
	c.mu.Lock()
	if c.cache.imageClient == nil {
		c.cache.imageClient = NewImageClient(c.client.Image)
	}
	defer c.mu.Unlock()
	return c.cache.imageClient
}

func (c *mockClient) ISO() ISOClient {
	c.mu.Lock()
	if c.cache.isoClient == nil {
		c.cache.isoClient = NewISOClient(c.client.ISO)
	}
	defer c.mu.Unlock()
	return c.cache.isoClient
}

func (c *mockClient) Location() LocationClient {
	c.mu.Lock()
	if c.cache.locationClient == nil {
		c.cache.locationClient = NewLocationClient(c.client.Location)
	}
	defer c.mu.Unlock()
	return c.cache.locationClient
}

func (c *mockClient) LoadBalancer() LoadBalancerClient {
	c.mu.Lock()
	if c.cache.loadBalancerClient == nil {
		c.cache.loadBalancerClient = NewLoadBalancerClient(c.client.LoadBalancer)
	}
	defer c.mu.Unlock()
	return c.cache.loadBalancerClient
}
func (c *mockClient) LoadBalancerType() LoadBalancerTypeClient {
	c.mu.Lock()
	if c.cache.loadBalancerTypeClient == nil {
		c.cache.loadBalancerTypeClient = NewLoadBalancerTypeClient(c.client.LoadBalancerType)
	}
	defer c.mu.Unlock()
	return c.cache.loadBalancerTypeClient
}
func (c *mockClient) Network() NetworkClient {
	c.mu.Lock()
	if c.cache.networkClient == nil {
		c.cache.networkClient = NewNetworkClient(c.client.Network)
	}
	defer c.mu.Unlock()
	return c.cache.networkClient
}

func (c *mockClient) Server() ServerClient {
	c.mu.Lock()
	if c.cache.serverClient == nil {
		c.cache.serverClient = NewServerClient(c.client.Server)
	}
	defer c.mu.Unlock()
	return c.cache.serverClient
}

func (c *mockClient) ServerType() ServerTypeClient {
	c.mu.Lock()
	if c.cache.serverTypeClient == nil {
		c.cache.serverTypeClient = NewServerTypeClient(c.client.ServerType)
	}
	defer c.mu.Unlock()
	return c.cache.serverTypeClient
}

func (c *mockClient) SSHKey() SSHKeyClient {
	c.mu.Lock()
	if c.cache.sshKeyClient == nil {
		c.cache.sshKeyClient = NewSSHKeyClient(c.client.SSHKey)
	}
	defer c.mu.Unlock()
	return c.cache.sshKeyClient
}
func (c *mockClient) RDNS() RDNSClient {
	c.mu.Lock()
	if c.cache.rdnsClient == nil {
		c.cache.rdnsClient = NewRDNSClient(c.client.RDNS)
	}
	defer c.mu.Unlock()
	return c.cache.rdnsClient
}

func (c *mockClient) Volume() VolumeClient {
	c.mu.Lock()
	if c.cache.volumeClient == nil {
		c.cache.volumeClient = NewVolumeClient(c.client.Volume)
	}
	defer c.mu.Unlock()
	return c.cache.volumeClient
}

func (c *mockClient) PlacementGroup() PlacementGroupClient {
	c.mu.Lock()
	if c.cache.placementGroupClient == nil {
		c.cache.placementGroupClient = NewPlacementGroupClient(c.client.PlacementGroup)
	}
	defer c.mu.Unlock()
	return c.cache.placementGroupClient
}
