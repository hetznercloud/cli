package mock

import (
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type Client struct {
	ActionClient           *MockActionClient
	CertificateClient      *MockCertificateClient
	DatacenterClient       *MockDatacenterClient
	FirewallClient         *MockFirewallClient
	FloatingIPClient       *MockFloatingIPClient
	PrimaryIPClient        *MockPrimaryIPClient
	ImageClient            *MockImageClient
	LocationClient         *MockLocationClient
	LoadBalancerClient     *MockLoadBalancerClient
	LoadBalancerTypeClient *MockLoadBalancerTypeClient
	NetworkClient          *MockNetworkClient
	ServerClient           *MockServerClient
	ServerTypeClient       *MockServerTypeClient
	SSHKeyClient           *MockSSHKeyClient
	VolumeClient           *MockVolumeClient
	ISOClient              *MockISOClient
	PlacementGroupClient   *MockPlacementGroupClient
	RDNSClient             *MockRDNSClient
	PricingClient          *MockPricingClient
	StorageBoxClient       *MockStorageBoxClient
	StorageBoxTypeClient   *MockStorageBoxTypeClient
	ZoneClient             *MockZoneClient
}

func NewMockClient(ctrl *gomock.Controller) *Client {

	primaryIPClient := NewMockPrimaryIPClient(ctrl)
	// We call .Names() in command constructors to provide completions. Since we do not want to mock each
	// completion function, we ignore the call to .Names() here.
	primaryIPClient.EXPECT().Names(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	return &Client{
		ActionClient:           NewMockActionClient(ctrl),
		CertificateClient:      NewMockCertificateClient(ctrl),
		DatacenterClient:       NewMockDatacenterClient(ctrl),
		FirewallClient:         NewMockFirewallClient(ctrl),
		FloatingIPClient:       NewMockFloatingIPClient(ctrl),
		PrimaryIPClient:        primaryIPClient,
		ImageClient:            NewMockImageClient(ctrl),
		LocationClient:         NewMockLocationClient(ctrl),
		LoadBalancerClient:     NewMockLoadBalancerClient(ctrl),
		LoadBalancerTypeClient: NewMockLoadBalancerTypeClient(ctrl),
		NetworkClient:          NewMockNetworkClient(ctrl),
		ServerClient:           NewMockServerClient(ctrl),
		ServerTypeClient:       NewMockServerTypeClient(ctrl),
		SSHKeyClient:           NewMockSSHKeyClient(ctrl),
		VolumeClient:           NewMockVolumeClient(ctrl),
		ISOClient:              NewMockISOClient(ctrl),
		PlacementGroupClient:   NewMockPlacementGroupClient(ctrl),
		RDNSClient:             NewMockRDNSClient(ctrl),
		PricingClient:          NewMockPricingClient(ctrl),
		StorageBoxClient:       NewMockStorageBoxClient(ctrl),
		StorageBoxTypeClient:   NewMockStorageBoxTypeClient(ctrl),
		ZoneClient:             NewMockZoneClient(ctrl),
	}
}

func (c *Client) Action() hcapi2.ActionClient {
	return c.ActionClient
}

func (c *Client) Certificate() hcapi2.CertificateClient {
	return c.CertificateClient
}

func (c *Client) Datacenter() hcapi2.DatacenterClient {
	return c.DatacenterClient
}

func (c *Client) Firewall() hcapi2.FirewallClient {
	return c.FirewallClient
}

func (c *Client) FloatingIP() hcapi2.FloatingIPClient {
	return c.FloatingIPClient
}

func (c *Client) PrimaryIP() hcapi2.PrimaryIPClient {
	return c.PrimaryIPClient
}

func (c *Client) Image() hcapi2.ImageClient {
	return c.ImageClient
}

func (c *Client) ISO() hcapi2.ISOClient {
	return c.ISOClient
}

func (c *Client) Location() hcapi2.LocationClient {
	return c.LocationClient
}

func (c *Client) LoadBalancer() hcapi2.LoadBalancerClient {
	return c.LoadBalancerClient
}

func (c *Client) LoadBalancerType() hcapi2.LoadBalancerTypeClient {
	return c.LoadBalancerTypeClient
}

func (c *Client) Network() hcapi2.NetworkClient {
	return c.NetworkClient
}

func (c *Client) Server() hcapi2.ServerClient {
	return c.ServerClient
}

func (c *Client) ServerType() hcapi2.ServerTypeClient {
	return c.ServerTypeClient
}

func (c *Client) SSHKey() hcapi2.SSHKeyClient {
	return c.SSHKeyClient
}

func (c *Client) Volume() hcapi2.VolumeClient {
	return c.VolumeClient
}
func (c *Client) RDNS() hcapi2.RDNSClient {
	return c.RDNSClient
}

func (c *Client) PlacementGroup() hcapi2.PlacementGroupClient {
	return c.PlacementGroupClient
}

func (c *Client) Pricing() hcapi2.PricingClient {
	return c.PricingClient
}

func (c *Client) StorageBox() hcapi2.StorageBoxClient {
	return c.StorageBoxClient
}

func (c *Client) StorageBoxType() hcapi2.StorageBoxTypeClient {
	return c.StorageBoxTypeClient
}

func (c *Client) Zone() hcapi2.ZoneClient {
	return c.ZoneClient
}

func (*Client) WithOpts(_ ...hcloud.ClientOption) {
	// no-op
}

func (*Client) FromConfig(_ config.Config) {
	// no-op
}
