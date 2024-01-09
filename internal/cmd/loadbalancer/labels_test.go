package loadbalancer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{ID: 123}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		Update(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, _, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.LoadBalancer{
			ID: 123,
			Labels: map[string]string{
				"key": "value",
			},
		}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		Update(gomock.Any(), &hcloud.LoadBalancer{ID: 123}, hcloud.LoadBalancerUpdateOpts{
			Labels: make(map[string]string),
		})

	out, _, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from Load Balancer 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
