package loadbalancer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestChangeAlgorithm(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ChangeAlgorithmCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		ChangeAlgorithm(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerChangeAlgorithmOpts{
			Type: hcloud.LoadBalancerAlgorithmTypeLeastConnections,
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"123", "--algorithm-type", "least_connections"})

	expOut := "Algorithm for Load Balancer 123 was changed\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
