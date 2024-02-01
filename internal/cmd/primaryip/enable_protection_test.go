package primaryip_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestEnableProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.EnableProtectionCmd.CobraCommand(fx.State())
	action := &hcloud.Action{ID: 1}
	ip := &hcloud.PrimaryIP{ID: 13}
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			ip,
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		ChangeProtection(
			gomock.Any(),
			hcloud.PrimaryIPChangeProtectionOpts{
				ID:     13,
				Delete: true,
			},
		).
		Return(
			action,
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), gomock.Any(), action).Return(nil)
	out, errOut, err := fx.Run(cmd, []string{"13"})

	expOut := "Resource protection enabled for primary IP 13\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
