package primaryip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestEnableProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := EnableProtectionCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer, fx.ActionWaiter)
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
				Delete: true,
			},
		).
		Return(
			action,
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), action).Return(nil)
	out, err := fx.Run(cmd, []string{"13"})

	expOut := "Resource protection enabled for primary IP 13\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
