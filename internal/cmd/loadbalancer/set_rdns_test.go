package loadbalancer

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestSetRDNS(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := SetRDNSCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	loadBalancer := &hcloud.LoadBalancer{
		ID: 123,
		PublicNet: hcloud.LoadBalancerPublicNet{
			IPv4: hcloud.LoadBalancerPublicNetIPv4{
				IP: net.ParseIP("192.168.2.1"),
			},
		},
	}

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(loadBalancer, nil, nil)
	fx.Client.RDNSClient.EXPECT().
		ChangeDNSPtr(gomock.Any(), loadBalancer, loadBalancer.PublicNet.IPv4.IP, hcloud.Ptr("example.com")).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"--hostname", "example.com", "test"})

	expOut := "Reverse DNS of Load Balancer test changed\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
