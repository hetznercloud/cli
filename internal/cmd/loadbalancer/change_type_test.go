package loadbalancer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestChangeType(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.ChangeTypeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	lbType := &hcloud.LoadBalancerType{ID: 321, Name: "lb21"}

	fx.Client.LoadBalancer.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerType.EXPECT().
		Get(gomock.Any(), "lb21").
		Return(lbType, nil, nil)
	fx.Client.LoadBalancer.EXPECT().
		ChangeType(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerChangeTypeOpts{
			LoadBalancerType: lbType,
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "lb21"})

	expOut := "LoadBalancer 123 changed to type lb21\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
