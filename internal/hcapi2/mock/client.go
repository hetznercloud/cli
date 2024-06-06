package hcapi2_mock

import (
	"github.com/golang/mock/gomock"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type MockClient struct {
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
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	return &MockClient{
		ActionClient:           NewMockActionClient(ctrl),
		CertificateClient:      NewMockCertificateClient(ctrl),
		DatacenterClient:       NewMockDatacenterClient(ctrl),
		FirewallClient:         NewMockFirewallClient(ctrl),
		FloatingIPClient:       NewMockFloatingIPClient(ctrl),
		PrimaryIPClient:        NewMockPrimaryIPClient(ctrl),
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
	}
}

func (c *MockClient) Action() hcapi2.ActionClient {
	return c.ActionClient
}

func (c *MockClient) Certificate() hcapi2.CertificateClient {
	return c.CertificateClient
}

func (c *MockClient) Datacenter() hcapi2.DatacenterClient {
	return c.DatacenterClient
}

func (c *MockClient) Firewall() hcapi2.FirewallClient {
	return c.FirewallClient
}

func (c *MockClient) FloatingIP() hcapi2.FloatingIPClient {
	return c.FloatingIPClient
}

func (c *MockClient) PrimaryIP() hcapi2.PrimaryIPClient {
	return c.PrimaryIPClient
}

func (c *MockClient) Image() hcapi2.ImageClient {
	return c.ImageClient
}

func (c *MockClient) ISO() hcapi2.ISOClient {
	return c.ISOClient
}

func (c *MockClient) Location() hcapi2.LocationClient {
	return c.LocationClient
}

func (c *MockClient) LoadBalancer() hcapi2.LoadBalancerClient {
	return c.LoadBalancerClient
}

func (c *MockClient) LoadBalancerType() hcapi2.LoadBalancerTypeClient {
	return c.LoadBalancerTypeClient
}

func (c *MockClient) Network() hcapi2.NetworkClient {
	return c.NetworkClient
}

func (c *MockClient) Server() hcapi2.ServerClient {
	return c.ServerClient
}

func (c *MockClient) ServerType() hcapi2.ServerTypeClient {
	return c.ServerTypeClient
}

func (c *MockClient) SSHKey() hcapi2.SSHKeyClient {
	return c.SSHKeyClient
}

func (c *MockClient) Volume() hcapi2.VolumeClient {
	return c.VolumeClient
}
func (c *MockClient) RDNS() hcapi2.RDNSClient {
	return c.RDNSClient
}

func (c *MockClient) PlacementGroup() hcapi2.PlacementGroupClient {
	return c.PlacementGroupClient
}

func (*MockClient) WithOpts(_ ...hcloud.ClientOption) {
	// no-op
}

func (*MockClient) FromConfig(_ config.Config) {
	// no-op
}
