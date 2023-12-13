package loadbalancer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAddService(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := AddServiceCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AddService(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerAddServiceOpts{
			Protocol:        hcloud.LoadBalancerServiceProtocolHTTP,
			ListenPort:      hcloud.Ptr(80),
			DestinationPort: hcloud.Ptr(8080),
			HTTP: &hcloud.LoadBalancerAddServiceOptsHTTP{
				StickySessions: hcloud.Ptr(false),
				RedirectHTTP:   hcloud.Ptr(false),
			},
			Proxyprotocol: hcloud.Ptr(false),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--protocol", "http", "--listen-port", "80", "--destination-port", "8080"})

	expOut := "Service was added to Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
