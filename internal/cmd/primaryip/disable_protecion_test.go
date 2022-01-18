package primaryip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestEnable(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DisableProtectionCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer, fx.ActionWaiter)
	action := &hcloud.Action{ID: 1}
	primaryip := &hcloud.PrimaryIP{ID: 13}
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			primaryip,
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		ChangeProtection(
			gomock.Any(),
			hcloud.PrimaryIPChangeProtectionOpts{
				ID:     13,
				Delete: false,
			},
		).
		Return(
			action,
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), action).Return(nil)
	out, err := fx.Run(cmd, []string{"13"})

	expOut := "Primary IP 13 protection disabled"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
