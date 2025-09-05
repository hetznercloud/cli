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
	StorageBox() StorageBoxClient
	StorageBoxType() StorageBoxTypeClient
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
	storageBoxClient       StorageBoxClient
	storageBoxTypeClient   StorageBoxTypeClient
}

type client struct {
	client *hcloud.Client
	cache  clientCache

	mu   sync.Mutex
	opts []hcloud.ClientOption
}

// NewClient creates a new CLI API client extending hcloud.Client.
func NewClient(opts ...hcloud.ClientOption) Client {
	c := &client{
		opts: opts,
	}
	c.update()
	return c
}

func (c *client) WithOpts(opts ...hcloud.ClientOption) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opts = append(c.opts, opts...)
	c.update()
}

func (c *client) update() {
	c.client = hcloud.NewClient(c.opts...)
	c.cache = clientCache{}
}

func (c *client) Action() ActionClient {
	c.mu.Lock()
	if c.cache.actionClient == nil {
		c.cache.actionClient = NewActionClient(&c.client.Action)
	}
	defer c.mu.Unlock()
	return c.cache.actionClient
}

func (c *client) Certificate() CertificateClient {
	c.mu.Lock()
	if c.cache.certificateClient == nil {
		c.cache.certificateClient = NewCertificateClient(&c.client.Certificate)
	}
	defer c.mu.Unlock()
	return c.cache.certificateClient
}

func (c *client) Datacenter() DatacenterClient {
	c.mu.Lock()
	if c.cache.datacenterClient == nil {
		c.cache.datacenterClient = NewDatacenterClient(&c.client.Datacenter)
	}
	defer c.mu.Unlock()
	return c.cache.datacenterClient
}

func (c *client) Firewall() FirewallClient {
	c.mu.Lock()
	if c.cache.firewallClient == nil {
		c.cache.firewallClient = NewFirewallClient(&c.client.Firewall)
	}
	defer c.mu.Unlock()
	return c.cache.firewallClient
}

func (c *client) FloatingIP() FloatingIPClient {
	c.mu.Lock()
	if c.cache.floatingIPClient == nil {
		c.cache.floatingIPClient = NewFloatingIPClient(&c.client.FloatingIP)
	}
	defer c.mu.Unlock()
	return c.cache.floatingIPClient
}

func (c *client) PrimaryIP() PrimaryIPClient {
	c.mu.Lock()
	if c.cache.primaryIPClient == nil {
		c.cache.primaryIPClient = NewPrimaryIPClient(&c.client.PrimaryIP)
	}
	defer c.mu.Unlock()
	return c.cache.primaryIPClient
}

func (c *client) Image() ImageClient {
	c.mu.Lock()
	if c.cache.imageClient == nil {
		c.cache.imageClient = NewImageClient(&c.client.Image)
	}
	defer c.mu.Unlock()
	return c.cache.imageClient
}

func (c *client) ISO() ISOClient {
	c.mu.Lock()
	if c.cache.isoClient == nil {
		c.cache.isoClient = NewISOClient(&c.client.ISO)
	}
	defer c.mu.Unlock()
	return c.cache.isoClient
}

func (c *client) Location() LocationClient {
	c.mu.Lock()
	if c.cache.locationClient == nil {
		c.cache.locationClient = NewLocationClient(&c.client.Location)
	}
	defer c.mu.Unlock()
	return c.cache.locationClient
}

func (c *client) LoadBalancer() LoadBalancerClient {
	c.mu.Lock()
	if c.cache.loadBalancerClient == nil {
		c.cache.loadBalancerClient = NewLoadBalancerClient(&c.client.LoadBalancer)
	}
	defer c.mu.Unlock()
	return c.cache.loadBalancerClient
}
func (c *client) LoadBalancerType() LoadBalancerTypeClient {
	c.mu.Lock()
	if c.cache.loadBalancerTypeClient == nil {
		c.cache.loadBalancerTypeClient = NewLoadBalancerTypeClient(&c.client.LoadBalancerType)
	}
	defer c.mu.Unlock()
	return c.cache.loadBalancerTypeClient
}
func (c *client) Network() NetworkClient {
	c.mu.Lock()
	if c.cache.networkClient == nil {
		c.cache.networkClient = NewNetworkClient(&c.client.Network)
	}
	defer c.mu.Unlock()
	return c.cache.networkClient
}

func (c *client) Server() ServerClient {
	c.mu.Lock()
	if c.cache.serverClient == nil {
		c.cache.serverClient = NewServerClient(&c.client.Server)
	}
	defer c.mu.Unlock()
	return c.cache.serverClient
}

func (c *client) ServerType() ServerTypeClient {
	c.mu.Lock()
	if c.cache.serverTypeClient == nil {
		c.cache.serverTypeClient = NewServerTypeClient(&c.client.ServerType)
	}
	defer c.mu.Unlock()
	return c.cache.serverTypeClient
}

func (c *client) SSHKey() SSHKeyClient {
	c.mu.Lock()
	if c.cache.sshKeyClient == nil {
		c.cache.sshKeyClient = NewSSHKeyClient(&c.client.SSHKey)
	}
	defer c.mu.Unlock()
	return c.cache.sshKeyClient
}
func (c *client) RDNS() RDNSClient {
	c.mu.Lock()
	if c.cache.rdnsClient == nil {
		c.cache.rdnsClient = NewRDNSClient(&c.client.RDNS)
	}
	defer c.mu.Unlock()
	return c.cache.rdnsClient
}

func (c *client) Volume() VolumeClient {
	c.mu.Lock()
	if c.cache.volumeClient == nil {
		c.cache.volumeClient = NewVolumeClient(&c.client.Volume)
	}
	defer c.mu.Unlock()
	return c.cache.volumeClient
}

func (c *client) PlacementGroup() PlacementGroupClient {
	c.mu.Lock()
	if c.cache.placementGroupClient == nil {
		c.cache.placementGroupClient = NewPlacementGroupClient(&c.client.PlacementGroup)
	}
	defer c.mu.Unlock()
	return c.cache.placementGroupClient
}

func (c *client) Pricing() PricingClient {
	c.mu.Lock()
	if c.cache.pricingClient == nil {
		c.cache.pricingClient = NewPricingClient(&c.client.Pricing)
	}
	defer c.mu.Unlock()
	return c.cache.pricingClient
}

func (c *client) StorageBox() StorageBoxClient {
	c.mu.Lock()
	if c.cache.storageBoxClient == nil {
		c.cache.storageBoxClient = NewStorageBoxClient(&c.client.StorageBox)
	}
	defer c.mu.Unlock()
	return c.cache.storageBoxClient
}

func (c *client) StorageBoxType() StorageBoxTypeClient {
	c.mu.Lock()
	if c.cache.storageBoxTypeClient == nil {
		c.cache.storageBoxTypeClient = NewStorageBoxTypeClient(&c.client.StorageBoxType)
	}
	defer c.mu.Unlock()
	return c.cache.storageBoxTypeClient
}
