package firewall

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.FirewallClient.EXPECT().
		Create(gomock.Any(), hcloud.FirewallCreateOpts{
			Name:   "test",
			Labels: make(map[string]string),
		}).
		Return(hcloud.FirewallCreateResult{
			Firewall: &hcloud.Firewall{
				ID:   123,
				Name: "test",
			},
			Actions: []*hcloud.Action{{ID: 321}},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 321}})

	out, err := fx.Run(cmd, []string{"--name", "test"})

	expOut := "Firewall 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
