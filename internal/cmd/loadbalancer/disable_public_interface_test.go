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

func TestDisablePublicInterface(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.DisablePublicInterfaceCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		DisablePublicInterface(gomock.Any(), &hcloud.LoadBalancer{ID: 123}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123"})

	expOut := "Public interface of Load Balancer 123 was disabled\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
