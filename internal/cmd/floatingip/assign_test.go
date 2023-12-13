package floatingip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAssign(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := AssignCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "my-ip").
		Return(&hcloud.FloatingIP{ID: 123}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(&hcloud.Server{ID: 456}, nil, nil)
	fx.Client.FloatingIPClient.EXPECT().
		Assign(gomock.Any(), &hcloud.FloatingIP{ID: 123}, &hcloud.Server{ID: 456}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"my-ip", "my-server"})

	expOut := "Floating IP 123 assigned to server 456\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
