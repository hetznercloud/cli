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

func TestAddTargetServer(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := AddTargetCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(&hcloud.Server{ID: 321}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AddServerTarget(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerAddServerTargetOpts{
			Server:       &hcloud.Server{ID: 321},
			UsePrivateIP: hcloud.Ptr(false),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--server", "my-server"})

	expOut := "Target added to Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestAddTargetLabelSelector(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := AddTargetCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AddLabelSelectorTarget(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerAddLabelSelectorTargetOpts{
			Selector:     "my-label",
			UsePrivateIP: hcloud.Ptr(false),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--label-selector", "my-label"})

	expOut := "Target added to Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestAddTargetIP(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := AddTargetCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AddIPTarget(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerAddIPTargetOpts{
			IP: net.ParseIP("192.168.2.1"),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--ip", "192.168.2.1"})

	expOut := "Target added to Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
