package primaryip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUnAssign(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := UnAssignCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer, fx.ActionWaiter)
	action := &hcloud.Action{ID: 1}
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			&hcloud.PrimaryIP{ID: 13},
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		Unassign(
			gomock.Any(),
			int64(13),
		).
		Return(
			action,
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), gomock.Any(), action).Return(nil)

	out, _, err := fx.Run(cmd, []string{"13"})

	expOut := "Primary IP 13 was unassigned successfully\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
