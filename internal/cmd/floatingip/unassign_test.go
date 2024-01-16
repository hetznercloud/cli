package floatingip_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUnassign(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.UnassignCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "my-ip").
		Return(&hcloud.FloatingIP{ID: 123}, nil, nil)
	fx.Client.FloatingIPClient.EXPECT().
		Unassign(gomock.Any(), &hcloud.FloatingIP{ID: 123}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"my-ip"})

	expOut := "Floating IP 123 unassigned\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
