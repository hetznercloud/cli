package loadbalancer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDeleteService(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteServiceCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		DeleteService(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, 80).
		Return(&hcloud.Action{ID: 123}, nil, nil)

	out, _, err := fx.Run(cmd, []string{"123", "--listen-port", "80"})

	expOut := "Service on port 80 deleted from Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
