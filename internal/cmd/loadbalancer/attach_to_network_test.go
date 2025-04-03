package loadbalancer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAttachToNetwork(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.AttachToNetworkCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "my-network").
		Return(&hcloud.Network{ID: 321}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AttachToNetwork(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerAttachToNetworkOpts{
			Network: &hcloud.Network{ID: 321},
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "--network", "my-network"})

	expOut := "Load Balancer 123 attached to Network 321\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
