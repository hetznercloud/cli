package loadbalancer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDisablePublicInterface(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DisablePublicInterfaceCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		DisablePublicInterface(gomock.Any(), &hcloud.LoadBalancer{ID: 123}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123"})

	expOut := "Public interface of Load Balancer 123 was disabled\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
