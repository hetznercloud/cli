package loadbalancer_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestSetRDNS(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.SetRDNSCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	loadBalancer := &hcloud.LoadBalancer{
		ID: 123,
		PublicNet: hcloud.LoadBalancerPublicNet{
			IPv4: hcloud.LoadBalancerPublicNetIPv4{
				IP: net.ParseIP("192.168.2.1"),
			},
		},
	}

	fx.Client.LoadBalancer.EXPECT().
		Get(gomock.Any(), "test").
		Return(loadBalancer, nil, nil)
	fx.Client.RDNS.EXPECT().
		ChangeDNSPtr(gomock.Any(), loadBalancer, loadBalancer.PublicNet.IPv4.IP, hcloud.Ptr("example.com")).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--hostname", "example.com", "test"})

	expOut := "Reverse DNS of Load Balancer test changed\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
