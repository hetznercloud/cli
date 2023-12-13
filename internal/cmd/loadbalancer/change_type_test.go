package loadbalancer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestChangeType(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ChangeTypeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	lbType := &hcloud.LoadBalancerType{ID: 321, Name: "lb21"}

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerTypeClient.EXPECT().
		Get(gomock.Any(), "lb21").
		Return(lbType, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		ChangeType(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerChangeTypeOpts{
			LoadBalancerType: lbType,
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "lb21"})

	expOut := "LoadBalancer 123 changed to type lb21\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
