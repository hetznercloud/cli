package hcapi2

import (
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Client makes all API clients accessible via a single interface.
type Client interface {
	Action() ActionClient
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
	LoadBalancerType() LoadBalancerTypeClient
	ISO() ISOClient
	PlacementGroup() PlacementGroupClient
	RDNS() RDNSClient
	PrimaryIP() PrimaryIPClient
	Pricing() PricingClient
	WithOpts(...hcloud.ClientOption)
}

type clientCache struct {
	actionClient           ActionClient
	certificateClient      CertificateClient
	datacenterClient       DatacenterClient
	serverClient           ServerClient
	serverTypeClient       ServerTypeClient
	locationClient         LocationClient
	loadBalancerClient     LoadBalancerClient
	loadBalancerTypeClient LoadBalancerTypeClient
	networkClient          NetworkClient
	firewallClient         FirewallClient
	floatingIPClient       FloatingIPClient
	imageClient            ImageClient
	isoClient              ISOClient
	sshKeyClient           SSHKeyClient
	volumeClient           VolumeClient
	placementGroupClient   PlacementGroupClient
	rdnsClient             RDNSClient
	primaryIPClient        PrimaryIPClient
	pricingClient          PricingClient
}

type ActualClient struct {
	Client *hcloud.Client
	cache  clientCache

	mu   sync.Mutex
	opts []hcloud.ClientOption
}

// NewClient creates a new CLI API ActualClient extending hcloud.Client.
func NewClient(opts ...hcloud.ClientOption) Client {
	c := &ActualClient{
		opts: opts,
	}
	c.update()
	return c
}

func (c *ActualClient) WithOpts(opts ...hcloud.ClientOption) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts = append(c.opts, opts...)
	c.update()
}

func (c *ActualClient) update() {
	c.Client = hcloud.NewClient(c.opts...)
	c.cache = clientCache{}
}

func (c *ActualClient) Action() ActionClient {
	c.mu.Lock()
	if c.cache.actionClient == nil {
		c.cache.actionClient = NewActionClient(&c.Client.Action)
	}
	defer c.mu.Unlock()
	return c.cache.actionClient
}

func (c *ActualClient) Certificate() CertificateClient {
	c.mu.Lock()
	if c.cache.certificateClient == nil {
		c.cache.certificateClient = NewCertificateClient(&c.Client.Certificate)
	}
	defer c.mu.Unlock()
	return c.cache.certificateClient
}

func (c *ActualClient) Datacenter() DatacenterClient {
	c.mu.Lock()
	if c.cache.datacenterClient == nil {
		c.cache.datacenterClient = NewDatacenterClient(&c.Client.Datacenter)
	}
	defer c.mu.Unlock()
	return c.cache.datacenterClient
}

func (c *ActualClient) Firewall() FirewallClient {
	c.mu.Lock()
	if c.cache.firewallClient == nil {
		c.cache.firewallClient = NewFirewallClient(&c.Client.Firewall)
	}
	defer c.mu.Unlock()
	return c.cache.firewallClient
}

func (c *ActualClient) FloatingIP() FloatingIPClient {
	c.mu.Lock()
	if c.cache.floatingIPClient == nil {
		c.cache.floatingIPClient = NewFloatingIPClient(&c.Client.FloatingIP)
	}
	defer c.mu.Unlock()
	return c.cache.floatingIPClient
}

func (c *ActualClient) PrimaryIP() PrimaryIPClient {
	c.mu.Lock()
	if c.cache.primaryIPClient == nil {
		c.cache.primaryIPClient = NewPrimaryIPClient(&c.Client.PrimaryIP)
	}
	defer c.mu.Unlock()
	return c.cache.primaryIPClient
}

func (c *ActualClient) Image() ImageClient {
	c.mu.Lock()
	if c.cache.imageClient == nil {
		c.cache.imageClient = NewImageClient(&c.Client.Image)
	}
	defer c.mu.Unlock()
	return c.cache.imageClient
}

func (c *ActualClient) ISO() ISOClient {
	c.mu.Lock()
	if c.cache.isoClient == nil {
		c.cache.isoClient = NewISOClient(&c.Client.ISO)
	}
	defer c.mu.Unlock()
	return c.cache.isoClient
}

func (c *ActualClient) Location() LocationClient {
	c.mu.Lock()
	if c.cache.locationClient == nil {
		c.cache.locationClient = NewLocationClient(&c.Client.Location)
	}
	defer c.mu.Unlock()
	return c.cache.locationClient
}

func (c *ActualClient) LoadBalancer() LoadBalancerClient {
	c.mu.Lock()
	if c.cache.loadBalancerClient == nil {
		c.cache.loadBalancerClient = NewLoadBalancerClient(&c.Client.LoadBalancer)
	}
	defer c.mu.Unlock()
	return c.cache.loadBalancerClient
}
func (c *ActualClient) LoadBalancerType() LoadBalancerTypeClient {
	c.mu.Lock()
	if c.cache.loadBalancerTypeClient == nil {
		c.cache.loadBalancerTypeClient = NewLoadBalancerTypeClient(&c.Client.LoadBalancerType)
	}
	defer c.mu.Unlock()
	return c.cache.loadBalancerTypeClient
}
func (c *ActualClient) Network() NetworkClient {
	c.mu.Lock()
	if c.cache.networkClient == nil {
		c.cache.networkClient = NewNetworkClient(&c.Client.Network)
	}
	defer c.mu.Unlock()
	return c.cache.networkClient
}

func (c *ActualClient) Server() ServerClient {
	c.mu.Lock()
	if c.cache.serverClient == nil {
		c.cache.serverClient = NewServerClient(&c.Client.Server)
	}
	defer c.mu.Unlock()
	return c.cache.serverClient
}

func (c *ActualClient) ServerType() ServerTypeClient {
	c.mu.Lock()
	if c.cache.serverTypeClient == nil {
		c.cache.serverTypeClient = NewServerTypeClient(&c.Client.ServerType)
	}
	defer c.mu.Unlock()
	return c.cache.serverTypeClient
}

func (c *ActualClient) SSHKey() SSHKeyClient {
	c.mu.Lock()
	if c.cache.sshKeyClient == nil {
		c.cache.sshKeyClient = NewSSHKeyClient(&c.Client.SSHKey)
	}
	defer c.mu.Unlock()
	return c.cache.sshKeyClient
}
func (c *ActualClient) RDNS() RDNSClient {
	c.mu.Lock()
	if c.cache.rdnsClient == nil {
		c.cache.rdnsClient = NewRDNSClient(&c.Client.RDNS)
	}
	defer c.mu.Unlock()
	return c.cache.rdnsClient
}

func (c *ActualClient) Volume() VolumeClient {
	c.mu.Lock()
	if c.cache.volumeClient == nil {
		c.cache.volumeClient = NewVolumeClient(&c.Client.Volume)
	}
	defer c.mu.Unlock()
	return c.cache.volumeClient
}

func (c *ActualClient) PlacementGroup() PlacementGroupClient {
	c.mu.Lock()
	if c.cache.placementGroupClient == nil {
		c.cache.placementGroupClient = NewPlacementGroupClient(&c.Client.PlacementGroup)
	}
	defer c.mu.Unlock()
	return c.cache.placementGroupClient
}

func (c *ActualClient) Pricing() PricingClient {
	c.mu.Lock()
	if c.cache.pricingClient == nil {
		c.cache.pricingClient = NewPricingClient(&c.Client.Pricing)
	}
	defer c.mu.Unlock()
	return c.cache.pricingClient
}
