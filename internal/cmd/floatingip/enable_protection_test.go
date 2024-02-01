package floatingip_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestEnableProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.EnableProtectionCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.FloatingIP{ID: 123}, nil, nil)
	fx.Client.FloatingIPClient.EXPECT().
		ChangeProtection(gomock.Any(), &hcloud.FloatingIP{ID: 123}, hcloud.FloatingIPChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"test", "delete"})

	expOut := "Resource protection enabled for floating IP 123\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
